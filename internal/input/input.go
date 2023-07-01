package input

type ContentInput interface {
	// It returns a channel with the words to be processed. the channel is closed when all words are read.
	// this function can be called concurrently.
	Words() <-chan string
	// It starts the reading process. It must be called before Words() is called.
	Prepare() error
}

type ReaderOptions struct {
	WordsChannelSize int
}

func DefaultFileReaderOptions() ReaderOptions {
	channelSize := 8
	return ReaderOptions{
		WordsChannelSize: channelSize,
	}
}
