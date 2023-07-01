package storage

import (
	"context"
)

type Storage interface {
	// It stores the word in the storage by hash. If the hash already exists, the word is appended to the list of words
	Store(ctx context.Context, hashKey string, word string) error
	// It returns the list of words for the given hash. If the hash does not exist, it returns an empty list.
	Get(ctx context.Context, hashKey string) ([]string, error)
	// It returns the channel of anagrams with coma seperated for all words in the storage. The channel will be
	// closed either when all the anagrams have been read or when an error occurs, signalled through the error channel.
	AllAnagrams(ctx context.Context) (<-chan []string, <-chan error)
}
