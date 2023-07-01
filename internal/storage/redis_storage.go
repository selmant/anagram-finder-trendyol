package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/selmant/anagram-finder-trendyol/internal"
	log "github.com/sirupsen/logrus"
)

type RedisAnagramStorage struct {
	redisClient *redis.Client
}

func NewRedisAnagramStorage(redisOptions redis.Options) *RedisAnagramStorage {
	return &RedisAnagramStorage{redisClient: redis.NewClient(&redisOptions)}
}

func (s *RedisAnagramStorage) Store(ctx context.Context, word string, letterMap internal.AnagramLetterMap) error {
	return s.redisClient.SAdd(ctx, letterMap.AnagramHash(), word).Err()
}

func (s *RedisAnagramStorage) AllAnagrams(ctx context.Context) (<-chan []string, <-chan error) {
	results := make(chan []string, 1)
	errors := make(chan error, 1)
	go func() {
		defer func() {
			log.Info("Closing anagram results and errors channels")
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
	return results, nil
}
