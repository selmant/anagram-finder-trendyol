package storage_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestLocalStorageIsAnagramStorage(t *testing.T) {
	assert.Implements(t, (*storage.Storage)(nil), new(storage.LocalStorage))
}

func TestRedisStorageIsAnagramStorage(t *testing.T) {
	assert.Implements(t, (*storage.Storage)(nil), new(storage.RedisStorage))
}

func TestStorageFactoryReturnsLocalStorage(t *testing.T) {
	factory := storage.UnifiedStorageFactory{}
	localCfg := config.Config{StorageType: config.StorageTypeLocal}
	assert.IsType(t, new(storage.LocalStorage), factory.CreateStorage(&localCfg))
}

func TestStorageFactoryReturnsRedisStorage(t *testing.T) {
	factory := storage.UnifiedStorageFactory{}
	redisCfg := config.Config{StorageType: config.StorageTypeRedis}
	assert.IsType(t, new(storage.RedisStorage), factory.CreateStorage(&redisCfg))
}
