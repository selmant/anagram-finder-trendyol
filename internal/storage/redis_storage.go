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

func (s *RedisStorage) AllAnagrams(ctx context.Context) (<-chan []string, <-chan error) {
	results := make(chan []string, 1)
	errors := make(chan error, 1)
	go func() {
		defer func() {
			log.Info("Closing results and errors channels")
			close(results)
			close(errors)
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
			errors <- err
			return
		}
		for _, cmd := range piped {
			err = cmd.Err()
			if err != nil {
				errors <- err
				return
			}
			anagrams := cmd.(*redis.StringSliceCmd).Val()
			results <- anagrams
		}
	}()
	return results, errors
}
