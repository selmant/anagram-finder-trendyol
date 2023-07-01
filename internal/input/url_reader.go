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
	linesChannel chan string
	options      ReaderOptions
}

const ten = 10 * time.Second // 10s

func NewURLReader(url string, options ReaderOptions) URLReader {
	linesChannel := make(chan string, options.WordsChannelSize)
	var netClient = &http.Client{
		Timeout: ten,
	}

	return URLReader{
		url:          url,
		linesChannel: linesChannel,
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
			f.linesChannel <- line
		}
		close(f.linesChannel)
		resp.Body.Close()
	}()

	return nil
}

func (f *URLReader) Lines(_ context.Context) <-chan string {
	out := make(chan string, 1)
	go func() {
		for line := range f.linesChannel {
			out <- line
		}
		close(out)
	}()
	return out
}
