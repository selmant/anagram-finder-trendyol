package storage_test

import (
	"context"
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestWordIsStoredInLocalStorage(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ls := storage.NewLocalStorage()
	err := ls.Store(ctx, "key", "test")
	assert.NoError(err)

	words, err := ls.Get(ctx, "key")
	assert.NoError(err)
	assert.Equal(len(words), 1)
	assert.Equal(words[0], "test")
}

func TestMultipleWordsAreStoredInLocalStorage(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ls := storage.NewLocalStorage()
	err := ls.Store(ctx, "key", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key", "test")
	assert.NoError(err)

	words, err := ls.Get(ctx, "key")
	assert.NoError(err)
	assert.Equal(len(words), 3)
	assert.Equal(words[0], "test")
}

func TestAnagramsReturnsAll(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ls := storage.NewLocalStorage()
	err := ls.Store(ctx, "key", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key2", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key3", "test")
	assert.NoError(err)

	all := ls.AllAnagrams(ctx)

	count := 0
	for r := range all {
		count++
		assert.NoError(r.Error)
	}
	assert.Equal(3, count)
}

func TestAnagramsReturnsAllWithMultipleKeys(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ls := storage.NewLocalStorage()
	err := ls.Store(ctx, "key", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key", "test2")
	assert.NoError(err)
	err = ls.Store(ctx, "key2", "test3")
	assert.NoError(err)

	all := ls.AllAnagrams(ctx)

	count := 0
	for r := range all {
		count++
		assert.NoError(r.Error)
	}
	assert.Equal(2, count)
}
