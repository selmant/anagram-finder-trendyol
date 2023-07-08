package storage_test

import (
	"context"
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LocalStorageTestSuite struct {
	suite.Suite
	ls  storage.LocalStorage
	ctx context.Context
}

func (suite *LocalStorageTestSuite) SetupTest() {
	suite.ls = *storage.NewLocalStorage()
	suite.ctx = context.Background()
}

func (suite *LocalStorageTestSuite) TestWordIsStoredInLocalStorage() {
	err := suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)

	words, err := suite.ls.Get(suite.ctx, "key")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), words, 1)
	assert.Equal(suite.T(), "test", words[0])
}

func (suite *LocalStorageTestSuite) TestMultipleWordsAreStoredInLocalStorage() {
	err := suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)

	words, err := suite.ls.Get(suite.ctx, "key")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), words, 3)
	assert.Equal(suite.T(), "test", words[0])
}

func (suite *LocalStorageTestSuite) TestAnagramsReturnsAll() {
	err := suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key2", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key3", "test")
	assert.NoError(suite.T(), err)

	all := suite.ls.AllAnagrams(suite.ctx)

	count := 0
	for r := range all {
		count++
		assert.NoError(suite.T(), r.Error)
	}
	assert.Equal(suite.T(), 3, count)
}

func (suite *LocalStorageTestSuite) TestAnagramsReturnsAllWithMultipleKeys() {
	err := suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key", "test2")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key2", "test3")
	assert.NoError(suite.T(), err)

	all := suite.ls.AllAnagrams(suite.ctx)

	count := 0
	for r := range all {
		count++
		assert.NoError(suite.T(), r.Error)
	}
	assert.Equal(suite.T(), 2, count)
}

func (suite *LocalStorageTestSuite) TestClearLocalStorage() {
	err := suite.ls.Store(suite.ctx, "key", "test")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key", "test2")
	assert.NoError(suite.T(), err)
	err = suite.ls.Store(suite.ctx, "key2", "test3")
	assert.NoError(suite.T(), err)

	words, err := suite.ls.Get(suite.ctx, "key2")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test3", words[0])

	err = suite.ls.Clear(suite.ctx)
	assert.NoError(suite.T(), err)

	words, err = suite.ls.Get(suite.ctx, "key")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), words, 0)

	words, err = suite.ls.Get(suite.ctx, "key2")
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), words)
}

func TestLocalStorageTestSuite(t *testing.T) {
	suite.Run(t, new(LocalStorageTestSuite))
}
