package storage

import (
	"context"

	"github.com/selmant/anagram-finder-trendyol/internal"
)

type AnagramStorage interface {
	// It stores the word in the storage by letter map.
	Store(ctx context.Context, word string, letterMap internal.AnagramLetterMap) error
	// It returns the channel of anagrams with coma seperated for all words in the storage. The channel will be
	// closed either when all the anagrams have been read or when an error occurs, signalled through the error channel.
	AllAnagrams(ctx context.Context) (<-chan []string, <-chan error)
}
