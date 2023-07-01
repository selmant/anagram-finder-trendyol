package input

import (
	"bufio"
	"context"
	"net/http"
	"time"
)

type URLReader struct {
	client       *http.Client
	url          string
	wordsChannel chan string
	options      ReaderOptions
}

const ten = 10 * time.Second // 10s

func NewURLReader(url string, options ReaderOptions) URLReader {
	wordsChannel := make(chan string, options.WordsChannelSize)
	var netClient = &http.Client{
		Timeout: ten,
	}

	return URLReader{
		url:          url,
		wordsChannel: wordsChannel,
		options:      options,
		client:       netClient,
	}
}

func (f *URLReader) Prepare(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return err
	}
	//nolint:bodyclose // body is closed in the goroutine
	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			f.wordsChannel <- line
		}
		close(f.wordsChannel)
		resp.Body.Close()
	}()

	return nil
}

func (f *URLReader) Words(_ context.Context) <-chan string {
	out := make(chan string, 1)
	go func() {
		for word := range f.wordsChannel {
			out <- word
		}
		close(out)
	}()
	return out
}
