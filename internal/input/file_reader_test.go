package input_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestFileReaderFileNotFound(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	fr := input.NewFileReader("unexistedfile.txt")

	err := fr.Prepare(ctx)
	assert.Error(err)
}

func TestFileReaderFileFound(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(f.Name(), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	fr := input.NewFileReader(f.Name())

	err = fr.Prepare(ctx)
	assert.NoError(err)
}

func TestFileReaderReadSingleLine(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(f.Name(), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	fr := input.NewFileReader(f.Name())

	err = fr.Prepare(ctx)
	assert.NoError(err)

	data := <-fr.Lines(ctx)
	assert.NoError(err)
	assert.Equal("test", data)
}

func TestFileReaderReadMultipleLines(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(f.Name(), []byte("test\ntest\ntest\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	fr := input.NewFileReader(f.Name())

	err = fr.Prepare(ctx)
	assert.NoError(err)

	count := 0
	for word := range fr.Lines(ctx) {
		count++
		assert.Equal("test", word)
	}

	assert.Equal(3, count)
}

func TestFileReaderConcurrentRead(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(f.Name(), []byte("test\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	fr := input.NewFileReader(f.Name())

	err = fr.Prepare(ctx)
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
