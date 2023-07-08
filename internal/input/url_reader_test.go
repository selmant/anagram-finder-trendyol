package input_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type URLReaderSuite struct {
	suite.Suite
	reader *input.URLReader
}

func (suite *URLReaderSuite) SetupTest() {
	suite.reader = nil
}

func (suite *URLReaderSuite) SetupSuite() {
	config.GlobalConfig = &config.Config{
		WordsChannelSize: 8,
		WorkerPoolSize:   16,
	}
}

func (suite *URLReaderSuite) TestURLNotFound() {
	ctx := context.Background()
	fr := input.NewURLReader("http://nonexisturl")
	err := fr.Prepare(ctx)
	assert.Error(suite.T(), err)
}

func (suite *URLReaderSuite) TestURLFound() {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test")
	}))
	defer ts.Close()
	suite.reader = input.NewURLReader(ts.URL)
	err := suite.reader.Prepare(ctx)
	assert.NoError(suite.T(), err)
}

func (suite *URLReaderSuite) TestReadMultipleLines() {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test\ntest\ntest")
	}))
	defer ts.Close()

	ur := input.NewURLReader(ts.URL)
	err := ur.Prepare(ctx)
	assert.NoError(suite.T(), err)

	count := 0
	for word := range ur.Lines(ctx) {
		count++
		assert.Equal(suite.T(), "test", word)
	}

	assert.Equal(suite.T(), 3, count)
}

func (suite *URLReaderSuite) TestConcurrentRead() {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest")
	}))
	defer ts.Close()

	fr := input.NewURLReader(ts.URL)

	err := fr.Prepare(ctx)
	assert.NoError(suite.T(), err)

	wg := sync.WaitGroup{}
	wg.Add(2)
	count1 := 0
	count2 := 0
	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count1++
			assert.Equal(suite.T(), "test", word)
		}
		assert.Greater(suite.T(), count1, 0)
		wg.Done()
	}()

	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count2++
			assert.Equal(suite.T(), "test", word)
		}
		assert.Greater(suite.T(), count2, 0)
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(suite.T(), 10, count1+count2)
}

func TestURLReaderSuite(t *testing.T) {
	suite.Run(t, new(URLReaderSuite))
}
