package input_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestFileReaderIsDataReader(t *testing.T) {
	assert.Implements(t, (*input.DataReader)(nil), new(input.FileReader))
}

func TestUrlReaderIsDataReader(t *testing.T) {
	assert.Implements(t, (*input.DataReader)(nil), new(input.URLReader))
}

func TestInputReaderFactoryReturnsFileReader(t *testing.T) {
	factory := input.UnifiedReaderFactory{}
	cfgFile := config.Config{}
	cfgFile.Input.File.Path = "test"
	reader, err := factory.CreateReader(&cfgFile)
	assert.Nil(t, err)
	assert.IsType(t, new(input.FileReader), reader)
}

func TestInputReaderFactoryReturnsUrlReader(t *testing.T) {
	factory := input.UnifiedReaderFactory{}
	cfgURL := config.Config{}
	cfgURL.Input.URL.URL = "test"
	reader, err := factory.CreateReader(&cfgURL)
	assert.Nil(t, err)
	assert.IsType(t, new(input.URLReader), reader)
}
