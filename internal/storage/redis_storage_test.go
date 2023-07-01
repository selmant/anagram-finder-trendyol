package storage_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/pkg/errors"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorageGet(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	client, mock := redismock.NewClientMock()
	redisStorage := storage.NewRedisStorage(*client)

	mock.ExpectSMembers("key").SetVal([]string{"test"})
	anagrams, err := redisStorage.Get(ctx, "key")
	assert.NoError(err)
	assert.Equal(len(anagrams), 1)
	assert.Equal(anagrams[0], "test")

	err = mock.ExpectationsWereMet()
	assert.NoError(err)
}

func TestWordIsStoredInRedisStorage(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	client, mock := redismock.NewClientMock()
	redisStorage := storage.NewRedisStorage(*client)
	mock.ExpectSAdd("key", "test").SetVal(1)

	err := redisStorage.Store(ctx, "key", "test")
	assert.NoError(err)

	err = mock.ExpectationsWereMet()
	assert.NoError(err)
}

func TestAllAnagrams(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	client, mock := redismock.NewClientMock()
	redisStorage := storage.NewRedisStorage(*client)

	mock.ExpectScan(0, "*", 0).SetVal([]string{"key", "key2"}, 1)
	mock.ExpectSMembers("key").SetVal([]string{"test", "sett"})
	mock.ExpectSMembers("key2").SetVal([]string{"asd", "dsa"})
	all, errs := redisStorage.AllAnagrams(ctx)

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
	err := mock.ExpectationsWereMet()
	assert.NoError(err)
	assert.Equal(count, 2)
}

func TestAllAnagramsWithError(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	client, mock := redismock.NewClientMock()
	redisStorage := storage.NewRedisStorage(*client)

	mock.ExpectScan(0, "*", 0).SetVal([]string{"key2"}, 1)
	mock.ExpectSMembers("key2").SetErr(errors.New("error"))
	all, errs := redisStorage.AllAnagrams(ctx)

	errCount := 0
	for all != nil || errs != nil {
		select {
		case _, ok := <-all:
			if !ok {
				all = nil
			}
		case _, ok := <-errs:
			if !ok {
				errs = nil
			} else {
				errCount++
			}
		}
	}
	err := mock.ExpectationsWereMet()
	assert.NoError(err)

	assert.Equal(errCount, 1)
}
