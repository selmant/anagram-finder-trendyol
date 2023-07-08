package storage_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisStorageTestSuite struct {
	suite.Suite
	client *redis.Client
	mock   redismock.ClientMock
	ctx    context.Context
}

func (suite *RedisStorageTestSuite) SetupTest() {
	suite.client, suite.mock = redismock.NewClientMock()
	suite.ctx = context.Background()
}

func (suite *RedisStorageTestSuite) TearDownTest() {
	err := suite.mock.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
	suite.mock.ClearExpect()
}

func (suite *RedisStorageTestSuite) TestRedisStorageGet() {
	assert := assert.New(suite.T())

	redisStorage := storage.NewRedisStorage(suite.client)

	suite.mock.ExpectSMembers("key").SetVal([]string{"test"})

	anagrams, err := redisStorage.Get(suite.ctx, "key")
	assert.NoError(err)
	assert.Len(anagrams, 1)
	assert.Equal("test", anagrams[0])
}

func (suite *RedisStorageTestSuite) TestWordIsStoredInRedisStorage() {
	assert := assert.New(suite.T())

	redisStorage := storage.NewRedisStorage(suite.client)
	suite.mock.ExpectSAdd("key", "test").SetVal(1)

	err := redisStorage.Store(suite.ctx, "key", "test")
	assert.NoError(err)
}

func (suite *RedisStorageTestSuite) TestAllAnagrams() {
	assert := assert.New(suite.T())

	redisStorage := storage.NewRedisStorage(suite.client)

	suite.mock.ExpectScan(0, "*", 0).SetVal([]string{"key", "key2"}, 1)
	suite.mock.ExpectSMembers("key").SetVal([]string{"test", "sett"})
	suite.mock.ExpectSMembers("key2").SetVal([]string{"asd", "dsa"})

	all := redisStorage.AllAnagrams(suite.ctx)

	count := 0
	for r := range all {
		count++
		assert.NoError(r.Error)
	}

	assert.Equal(2, count)
}

func (suite *RedisStorageTestSuite) TestAllAnagramsWithError() {
	assert := assert.New(suite.T())

	redisStorage := storage.NewRedisStorage(suite.client)

	suite.mock.ExpectScan(0, "*", 0).SetVal([]string{"key2"}, 1)
	suite.mock.ExpectSMembers("key2").SetErr(errors.New("error"))

	all := redisStorage.AllAnagrams(suite.ctx)

	errCount := 0
	for r := range all {
		errCount++
		assert.Error(r.Error)
	}

	assert.Equal(1, errCount)
}

func (suite *RedisStorageTestSuite) TestNewRedisClient() {
	assert := assert.New(suite.T())

	client := storage.NewRedisClient("localhost", 6379, "", 0)
	assert.NotNil(client)
}

func (suite *RedisStorageTestSuite) TestRedisClientClear() {
	redisStorage := storage.NewRedisStorage(suite.client)

	suite.mock.ExpectFlushDB()
	_ = redisStorage.Clear(suite.ctx)
	// TODO: check error
	// assert.NoError(suite.T(), err)
}

func TestRedisStorageTestSuite(t *testing.T) {
	suite.Run(t, new(RedisStorageTestSuite))
}
