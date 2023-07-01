package input_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestUrlReaderUrlNotFound(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	fr := input.NewURLReader("http://nonexisturl", input.DefaultFileReaderOptions())

	err := fr.Prepare(ctx)
	assert.Error(err)
}

func TestUrlReaderUrlFound(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test")
	}))
	defer ts.Close()

	fr := input.NewURLReader(ts.URL, input.DefaultFileReaderOptions())

	err := fr.Prepare(ctx)
	assert.NoError(err)
}

func TestUrlReaderReadMultipleLines(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test\ntest\ntest")
	}))
	defer ts.Close()

	ur := input.NewURLReader(ts.URL, input.DefaultFileReaderOptions())

	err := ur.Prepare(ctx)
	assert.NoError(err)

	count := 0
	for word := range ur.Lines(ctx) {
		count++
		assert.Equal("test", word)
	}

	assert.Equal(3, count)
}

func TestUrlReaderConcurrentRead(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest")
	}))
	defer ts.Close()

	fr := input.NewURLReader(ts.URL, input.DefaultFileReaderOptions())

	err := fr.Prepare(ctx)
	assert.NoError(err)

	wg := sync.WaitGroup{}
	wg.Add(2)
	count1 := 0
	count2 := 0
	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count1++
			assert.Equal("test", word)
		}
		assert.Greater(count1, 0)
		wg.Done()
	}()

	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count2++
			assert.Equal("test", word)
		}
		assert.Greater(count2, 0)
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(10, count1+count2)
}
