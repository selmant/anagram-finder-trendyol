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
	err = ls.Store(ctx, "key", "test")
	assert.NoError(err)
	err = ls.Store(ctx, "key", "test")
	assert.NoError(err)

	all, errs := ls.AllAnagrams(ctx)

	completed := false
	for !completed {
		select {
		case anagrams := <-all:
			assert.Equal(len(anagrams), 3)
		case err = <-errs:
			assert.NoError(err)
		default:
			completed = true
		}
	}
	assert.True(true)
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

	all, errs := ls.AllAnagrams(ctx)

	count := 0
	for all != nil || errs != nil {
		select {
		case _, ok := <-all:
			if !ok {
				all = nil
			} else {
				count++
			}
		case chanErr, ok := <-errs:
			if !ok {
				errs = nil
			} else {
				assert.NoError(chanErr)
			}
		}
	}
	assert.Equal(2, count)
}
