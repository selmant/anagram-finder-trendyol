package input_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InputReaderFactorySuite struct {
	suite.Suite
	factory input.UnifiedReaderFactory
}

func (suite *InputReaderFactorySuite) SetupTest() {
	suite.factory = input.UnifiedReaderFactory{}
}

func (suite *InputReaderFactorySuite) TestFileReaderIsDataReader() {
	assert.Implements(suite.T(), (*input.DataReader)(nil), new(input.FileReader))
}

func (suite *InputReaderFactorySuite) TestUrlReaderIsDataReader() {
	assert.Implements(suite.T(), (*input.DataReader)(nil), new(input.URLReader))
}

func (suite *InputReaderFactorySuite) TestInputReaderFactoryReturnsFileReader() {
	cfg := config.Config{}
	cfg.Input.File.Path = "test"
	reader, err := suite.factory.CreateReader(&cfg)
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), new(input.FileReader), reader)
}

func (suite *InputReaderFactorySuite) TestInputReaderFactoryReturnsUrlReader() {
	cfg := config.Config{}
	cfg.Input.URL.URL = "http://test"
	reader, err := suite.factory.CreateReader(&cfg)
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), new(input.URLReader), reader)
}

func TestInputReaderFactorySuite(t *testing.T) {
	suite.Run(t, new(InputReaderFactorySuite))
}
