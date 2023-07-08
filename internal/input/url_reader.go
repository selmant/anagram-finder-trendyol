package input

import (
	"bufio"
	"context"
	"net/http"
	"time"

	"github.com/selmant/anagram-finder-trendyol/app/config"
)

type URLReader struct {
	client       *http.Client
	url          string
	linesChannel chan string
}

const ten = 10 * time.Second // 10s

func NewURLReader(url string) *URLReader {
	linesChannel := make(chan string, config.GlobalConfig.WordsChannelSize)
	var netClient = &http.Client{
		Timeout: ten,
	}

	return &URLReader{
		url:          url,
		linesChannel: linesChannel,
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
		defer resp.Body.Close()
		defer close(f.linesChannel)

		for scanner.Scan() {
			line := scanner.Text()
			f.linesChannel <- line
		}
	}()

	return nil
}

func (f *URLReader) Lines(_ context.Context) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer close(out)

		for line := range f.linesChannel {
			out <- line
		}
	}()
	return out
}
