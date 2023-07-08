package input

import (
	"bufio"
	"context"
	"os"

	"github.com/selmant/anagram-finder-trendyol/app/config"
)

type FileReader struct {
	path         string
	linesChannel chan string
}

func NewFileReader(path string) *FileReader {
	linesChannel := make(chan string, config.GlobalConfig.WordsChannelSize)

	return &FileReader{
		path:         path,
		linesChannel: linesChannel,
	}
}

func (f *FileReader) Prepare(_ context.Context) error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	go func() {
		defer file.Close()
		defer close(f.linesChannel)

		for scanner.Scan() {
			line := scanner.Text()
			f.linesChannel <- line
		}
	}()

	return nil
}

func (f *FileReader) Lines(_ context.Context) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer close(out)

		for line := range f.linesChannel {
			out <- line
		}
	}()
	return out
}
