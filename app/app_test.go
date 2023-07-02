package app_test

import (
	"context"

	"github.com/stretchr/testify/assert"
)

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
