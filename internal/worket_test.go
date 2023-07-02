package internal_test

import (
	"context"
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/stretchr/testify/assert"
)

type MockJob struct {
	ProcessCallCount int
}

func (j *MockJob) Process(_ context.Context) error {
	j.ProcessCallCount++
	return nil
}

type MockJobWithError struct {
	ProcessCallCount int
}

func (j *MockJobWithError) Process(_ context.Context) error {
	j.ProcessCallCount++
	return assert.AnError
}

type MockStorage struct {
	StoreCallCount int
}

func (s *MockStorage) Store(_ context.Context, _, _ string) error {
	s.StoreCallCount++
	return nil
}

func (s *MockStorage) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	return nil, nil
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
	return nil, nil
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

func TestWorkePoolStart(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	workerCount := 2
	job := &MockJob{}
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.NoError(err)
	assert.Equal(workerCount, job.ProcessCallCount)
}

func TestWorkePoolStartWithError(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	workerCount := 2
	job := &MockJobWithError{}
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.Error(err)
	assert.Equal(workerCount, job.ProcessCallCount)
}

func TestReadAndMatchAnagram(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	mockStorage := &MockStorage{}
	mockRedaer := &MockInputReader{}
	mockRedaer.LinesChannel = make(chan string)
	go func() {
		mockRedaer.LinesChannel <- "abc"
		mockRedaer.LinesChannel <- "def"
		mockRedaer.LinesChannel <- "bac"
		mockRedaer.LinesChannel <- "fed"
		close(mockRedaer.LinesChannel)
	}()
	workerCount := 3

	job := internal.NewReadAndMatchAnagramJob(mockStorage, mockRedaer)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.NoError(err)
	assert.Equal(4, mockStorage.StoreCallCount)
	assert.Equal(workerCount, mockRedaer.LinesCallCount)
}

func TestReadAndMatchAnagramWithError(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	mockStorage := &MockStorage{}
	mockRedaer := &MockInputReader{}
	mockRedaer.LinesChannel = make(chan string)
	go func() {
		mockRedaer.LinesChannel <- "abc"
		mockRedaer.LinesChannel <- "def"
		mockRedaer.LinesChannel <- "bac"
		mockRedaer.LinesChannel <- "...123sdas"
		close(mockRedaer.LinesChannel)
	}()
	workerCount := 3

	job := internal.NewReadAndMatchAnagramJob(mockStorage, mockRedaer)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.Error(err)
	assert.Equal(3, mockStorage.StoreCallCount)
	assert.Equal(workerCount, mockRedaer.LinesCallCount)
}

func TestReadAndMatchAnagramWithStorageError(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	mockStorage := &MockStorageWithError{}
	mockRedaer := &MockInputReader{}
	mockRedaer.LinesChannel = make(chan string)
	go func() {
		mockRedaer.LinesChannel <- "abc"
		mockRedaer.LinesChannel <- "def"
		mockRedaer.LinesChannel <- "bac"
		mockRedaer.LinesChannel <- "fed"
		close(mockRedaer.LinesChannel)
	}()
	workerCount := 3

	job := internal.NewReadAndMatchAnagramJob(mockStorage, mockRedaer)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.Error(err)
	assert.Equal(4, mockStorage.StoreCallCount)
	assert.Equal(workerCount, mockRedaer.LinesCallCount)
}
