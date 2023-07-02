package storage

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type LocalStorage struct {
	storage map[string][]string
	lock    sync.RWMutex
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		storage: make(map[string][]string),
		lock:    sync.RWMutex{},
	}
}

func (s *LocalStorage) Store(_ context.Context, hashKey string, word string) error {
	s.lock.Lock()
	s.storage[hashKey] = append(s.storage[hashKey], word)
	s.lock.Unlock()
	return nil
}

func (s *LocalStorage) Get(_ context.Context, hashKey string) ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.storage[hashKey], nil
}

func (s *LocalStorage) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	out := make(chan []string, 1)
	errors := make(chan error, 1)
	go func() {
		for _, anagrams := range s.storage {
			out <- anagrams
		}
		close(out)
		close(errors)
		log.Info("All anagrams sent and channels are closed")
	}()
	return out, nil
}
