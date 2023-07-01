package storage

import (
	"context"
)

type LocalStorage struct {
	storage map[string][]string
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{storage: make(map[string][]string)}
}

func (s *LocalStorage) Store(_ context.Context, hashKey string, word string) error {
	s.storage[hashKey] = append(s.storage[hashKey], word)
	return nil
}

func (s *LocalStorage) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	out := make(chan []string, 1)
	go func() {
		for _, anagrams := range s.storage {
			out <- anagrams
		}
		close(out)
	}()
	return out, nil
}
