package app_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/stretchr/testify/assert"
)

//nolint:gochecknoinits // this is test file
func init() {
	config.GlobalConfig = &config.Config{}
	config.GlobalConfig.WorkerPoolSize = 1
	config.GlobalConfig.WordsChannelSize = 8
}

type MockStorage struct {
	StoreCallCount int
}

func (s *MockStorage) Store(_ context.Context, _, _ string) error {
	s.StoreCallCount++
	return nil
}

func (s *MockStorage) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	anagrams := make(chan []string)
	errs := make(chan error)
	go func() {
		anagrams <- []string{"test"}
		anagrams <- []string{"abc", "acb"}
		anagrams <- []string{"bac", "bca"}
		close(anagrams)
		close(errs)
	}()
	return anagrams, errs
}

func (s *MockStorage) Get(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

type MockStorageWithError struct {
	StoreCallCount int
}

func (s *MockStorageWithError) Store(_ context.Context, _, _ string) error {
	s.StoreCallCount++
	return assert.AnError
}

func (s *MockStorageWithError) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	errs := make(chan error)
	anagrams := make(chan []string)
	go func() {
		errs <- assert.AnError
		close(errs)
		close(anagrams)
	}()
	return anagrams, errs
}

func (s *MockStorageWithError) Get(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

type MockInputReader struct {
	LinesCallCount   int
	PrepareCallCount int
	LinesChannel     chan string
}

func (i *MockInputReader) Lines(_ context.Context) <-chan string {
	i.LinesCallCount++
	return i.LinesChannel
}
func (i *MockInputReader) Prepare(_ context.Context) error {
	i.PrepareCallCount++
	return nil
}

func TestAppPrintAnagrams(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockReader := &MockInputReader{}
	mockStorage := &MockStorage{}
	app := app.AnagramApplication{
		Input:          mockReader,
		AnagramStorage: mockStorage,
	}
	go func() {
		mockReader.LinesChannel = make(chan string)
		mockReader.LinesChannel <- "test"
		mockReader.LinesChannel <- "tets"
		close(mockReader.LinesChannel)
	}()
	err := app.PrintAnagrams(context.Background())
	assert.NoError(t, err)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, "abc, acb\nbac, bca\n", string(out))
}

func TestAppPrintAnagramsError(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockReader := &MockInputReader{}
	mockStorage := &MockStorageWithError{}
	app := app.AnagramApplication{
		Input:          mockReader,
		AnagramStorage: mockStorage,
	}
	go func() {
		mockReader.LinesChannel = make(chan string)
		mockReader.LinesChannel <- "test"
		mockReader.LinesChannel <- "tets"
		close(mockReader.LinesChannel)
	}()
	err := app.PrintAnagrams(context.Background())
	assert.Error(t, err)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, "", string(out))
}
