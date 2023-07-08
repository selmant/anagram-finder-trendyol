package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	redisClient.FlushDB(context.Background())
	return &RedisStorage{redisClient: redisClient}
}

func NewRedisClient(host string, port int, password string, db int) *redis.Client {
	var addr string
	if port != 0 {
		addr = fmt.Sprintf("%s:%d", host, port)
	} else {
		addr = host
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func (s *RedisStorage) Store(ctx context.Context, hashKey string, word string) error {
	return s.redisClient.SAdd(ctx, hashKey, word).Err()
}

func (s *RedisStorage) Get(ctx context.Context, hashKey string) ([]string, error) {
	return s.redisClient.SMembers(ctx, hashKey).Result()
}

func (s *RedisStorage) AllAnagrams(ctx context.Context) <-chan AnagramResult {
	results := make(chan AnagramResult, 1)
	go func() {
		defer func() {
			log.Info("Closing result channel")
			close(results)
		}()
		keys := s.redisClient.Scan(ctx, 0, "*", 0).Iterator()
		piped, err := s.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for keys.Next(ctx) {
				pipe.SMembers(ctx, keys.Val())
				log.Debugf("Anagrams for %s pipelined", keys.Val())
			}
			return nil
		})
		if err != nil {
			results <- AnagramResult{Error: err}
			return
		}
		for _, cmd := range piped {
			err = cmd.Err()
			if err != nil {
				results <- AnagramResult{cmd.String(), nil, err}
				return
			}
			anagrams := cmd.(*redis.StringSliceCmd).Val()
			results <- AnagramResult{cmd.String(), anagrams, nil}
		}
	}()
	return results
}

func (s *RedisStorage) Clear(ctx context.Context) error {
	return s.redisClient.FlushDB(ctx).Err()
}
