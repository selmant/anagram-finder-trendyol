package storage

import (
	"context"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	log "github.com/sirupsen/logrus"
)

type Storage interface {
	// It stores the word in the storage by hash. If the hash already exists, the word is appended to the list of words
	Store(ctx context.Context, hashKey string, word string) error
	// It returns the list of words for the given hash. If the hash does not exist, it returns an empty list.
	Get(ctx context.Context, hashKey string) ([]string, error)
	// It returns the channel of anagrams with coma seperated for all words in the storage. The channel will be
	// closed either when all the anagrams have been read or when an error occurs, signalled through the error channel.
	AllAnagrams(ctx context.Context) (<-chan []string, <-chan error)
}

type Factory interface {
	// It creates a new storage based on the configuration. It panics if the configuration is invalid.
	CreateStorage(cfg *config.Config) Storage
}

type UnifiedStorageFactory struct{}

func (f *UnifiedStorageFactory) CreateStorage(cfg *config.Config) Storage {
	switch cfg.StorageType {
	case config.StorageTypeLocal:
		return NewLocalStorage()
	case config.StorageTypeRedis:
		redisClient := NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
		if cmd := redisClient.Ping(context.Background()); cmd.Err() != nil {
			log.Fatal(cmd.Err())
		}
		return NewRedisStorage(redisClient)
	default:
		log.Fatalf("Invalid storage type: %s", cfg.StorageType)
		return nil
	}
}
