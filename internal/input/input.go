package input

import (
	"context"
	"log"

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
	CreateReader(cfg *config.Config) DataReader
}

type UnifiedReaderFactory struct{}

func (f *UnifiedReaderFactory) CreateReader(cfg *config.Config) DataReader {
	if cfg.Input.File.Path != "" {
		return NewFileReader(cfg.Input.File.Path)
	}
	if cfg.Input.URL.URL != "" {
		return NewURLReader(cfg.Input.URL.URL)
	}
	log.Fatal("Invalid input type")
	return nil
}
