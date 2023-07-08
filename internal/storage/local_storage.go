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

func (s *LocalStorage) AllAnagrams(_ context.Context) <-chan AnagramResult {
	out := make(chan AnagramResult, 1)
	go func() {
		for hashKey, anagrams := range s.storage {
			out <- AnagramResult{hashKey, anagrams, nil}
		}
		close(out)
		log.Info("All anagrams sent and channels are closed")
	}()
	return out
}

func (s *LocalStorage) Clear(_ context.Context) error {
	s.lock.Lock()
	s.storage = make(map[string][]string)
	s.lock.Unlock()
	return nil
}
