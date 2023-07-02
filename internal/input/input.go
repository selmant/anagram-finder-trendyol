package input

import "context"

type DataReader interface {
	// It returns a channel with the words to be processed. the channel is closed when all words are read.
	// this function can be called concurrently.
	Lines(ctx context.Context) <-chan string
	// It starts the reading process. It must be called before Words() is called.
	Prepare(ctx context.Context) error
}
