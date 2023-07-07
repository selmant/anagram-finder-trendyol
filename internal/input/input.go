package input

import (
	"context"
	"errors"

	"github.com/selmant/anagram-finder-trendyol/app/config"
)

type DataReader interface {
	// It returns a channel with the words to be processed. the channel is closed when all words are read.
	// this function can be called concurrently.
	Lines(ctx context.Context) <-chan string
	// It starts the reading process. It must be called before Words() is called.
	Prepare(ctx context.Context) error
}

type Factory interface {
	CreateReader(cfg *config.Config) (DataReader, error)
}

type UnifiedReaderFactory struct{}

func (f *UnifiedReaderFactory) CreateReader(cfg *config.Config) (DataReader, error) {
	if cfg.Input.File.Path != "" {
		return NewFileReader(cfg.Input.File.Path), nil
	}
	if cfg.Input.URL.URL != "" {
		return NewURLReader(cfg.Input.URL.URL), nil
	}
	return nil, errors.New("invalid input type")
}
