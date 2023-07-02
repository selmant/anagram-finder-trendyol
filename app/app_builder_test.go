package app_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
)

type MockReaderFactory struct{}

func (m MockReaderFactory) CreateReader(_ *config.Config) input.DataReader {
	return &MockInputReader{}
}

type MockStorageFactory struct{}

func (m MockStorageFactory) CreateStorage(_ *config.Config) storage.Storage {
	return &MockStorage{}
}

func TestBuilderGivesErrorWhenInputReaderIsNotSet(t *testing.T) {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithStorageFactory(&MockStorageFactory{})
	_, err := builder.Build()
	assert.Error(t, err)
}

func TestBuilderGivesErrorWhenStorageFactoryIsNotSet(t *testing.T) {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithReaderFactory(&MockReaderFactory{})
	_, err := builder.Build()
	assert.Error(t, err)
}

func TestBuilderGivesErrorWhenConfigIsNotSet(t *testing.T) {
	builder := app.NewAnagramApplicationBuilder().
		WithStorageFactory(&MockStorageFactory{}).
		WithReaderFactory(&MockReaderFactory{})
	_, err := builder.Build()
	assert.Error(t, err)
}

func TestAllDependenciesAreSet(t *testing.T) {
	builder := app.NewAnagramApplicationBuilder().
		WithConfig(&config.Config{}).
		WithStorageFactory(&MockStorageFactory{}).
		WithReaderFactory(&MockReaderFactory{})
	_, err := builder.Build()
	assert.NoError(t, err)
}
