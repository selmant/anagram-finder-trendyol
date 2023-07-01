package storage_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestLocalStorageIsAnagramStorage(t *testing.T) {
	assert.Implements(t, (*storage.AnagramStorage)(nil), new(storage.AnagramLocalStorage))
}

func TestRedisStorageIsAnagramStorage(t *testing.T) {
	assert.Implements(t, (*storage.AnagramStorage)(nil), new(storage.RedisAnagramStorage))
}
