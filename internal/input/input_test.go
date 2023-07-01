package input_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestFileReaderIsDataReader(t *testing.T) {
	assert.Implements(t, (*input.DataReader)(nil), new(input.FileReader))
}

func TestUrlReaderIsDataReader(t *testing.T) {
	assert.Implements(t, (*input.DataReader)(nil), new(input.URLReader))
}
