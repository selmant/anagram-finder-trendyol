package storage

import (
	"context"

	"github.com/selmant/anagram-finder-trendyol/internal"
)

type AnagramLocalStorage struct {
	storage map[string][]string
}

func NewAnagramLocalStorage() *AnagramLocalStorage {
	return &AnagramLocalStorage{storage: make(map[string][]string)}
}

func (s *AnagramLocalStorage) Store(_ context.Context, word string, letterMap internal.AnagramLetterMap) error {
	hash := letterMap.AnagramHash()
	s.storage[hash] = append(s.storage[hash], word)
	return nil
}

func (s *AnagramLocalStorage) AllAnagrams(_ context.Context) (<-chan []string, <-chan error) {
	out := make(chan []string, 1)
	go func() {
		for _, anagrams := range s.storage {
			out <- anagrams
		}
		close(out)
	}()
	return out, nil
}
