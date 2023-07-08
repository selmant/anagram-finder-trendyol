package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_input "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/input"
	mock_storage "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/storage"
	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	mockReaderFactory  *mock_input.MockFactory
	mockStorageFactory *mock_storage.MockFactory
	mockInputReader    *mock_input.MockDataReader
	mockStorage        *mock_storage.MockStorage
}

func (suite *BuilderTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockReaderFactory = mock_input.NewMockFactory(suite.ctrl)
	suite.mockStorageFactory = mock_storage.NewMockFactory(suite.ctrl)
	suite.mockInputReader = mock_input.NewMockDataReader(suite.ctrl)
	suite.mockStorage = mock_storage.NewMockStorage(suite.ctrl)
}

func (suite *BuilderTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *BuilderTestSuite) TestInputCfgIsEmpty() {
	builder := app.NewAnagramApplicationBuilder().
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestStorageFactoryIsEmpty() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithReaderFactory(suite.mockReaderFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestReaderFactoryIsEmpty() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithStorageFactory(suite.mockStorageFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestBuilderGivesErrorWhenInputReaderIsNotSet() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithStorageFactory(suite.mockStorageFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestBuilderGivesErrorWhenStorageFactoryIsNotSet() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithReaderFactory(suite.mockReaderFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestBuilderGivesErrorWhenConfigIsNotSet() {
	builder := app.NewAnagramApplicationBuilder().
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestInputInvalidStorage() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{Input: struct {
			File struct{ Path string }
			URL  struct{ URL string }
		}{}}).
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	suite.mockStorageFactory.EXPECT().CreateStorage(gomock.Any()).Return(nil, assert.AnError)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestNilStorage() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{Input: struct {
			File struct{ Path string }
			URL  struct{ URL string }
		}{}}).
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	suite.mockReaderFactory.EXPECT().CreateReader(gomock.Any()).Return(suite.mockInputReader, nil)
	suite.mockStorageFactory.EXPECT().CreateStorage(gomock.Any()).Return(nil, nil)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestNilInputReader() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{Input: struct {
			File struct{ Path string }
			URL  struct{ URL string }
		}{}}).
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	suite.mockReaderFactory.EXPECT().CreateReader(gomock.Any()).Return(nil, nil)
	suite.mockStorageFactory.EXPECT().CreateStorage(gomock.Any()).Return(suite.mockStorage, nil)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestInputInvalidReader() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{Input: struct {
			File struct{ Path string }
			URL  struct{ URL string }
		}{}}).
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	suite.mockReaderFactory.EXPECT().CreateReader(gomock.Any()).Return(nil, assert.AnError)
	suite.mockStorageFactory.EXPECT().CreateStorage(gomock.Any()).Return(suite.mockStorage, nil)
	_, err := builder.Build()
	assert.Error(suite.T(), err)
}

func (suite *BuilderTestSuite) TestAllDependenciesAreSet() {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithStorageFactory(suite.mockStorageFactory).
		WithReaderFactory(suite.mockReaderFactory)
	suite.mockReaderFactory.EXPECT().CreateReader(gomock.Any()).Return(suite.mockInputReader, nil)
	suite.mockStorageFactory.EXPECT().CreateStorage(gomock.Any()).Return(suite.mockStorage, nil)
	_, err := builder.Build()
	assert.NoError(suite.T(), err)
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
