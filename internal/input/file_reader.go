package input

import (
	"bufio"
	"context"
	"os"
)

type FileReader struct {
	path         string
	wordsChannel chan string
	options      ReaderOptions
}

func NewFileReader(path string, options ReaderOptions) FileReader {
	wordsChannel := make(chan string, options.WordsChannelSize)

	return FileReader{
		path:         path,
		wordsChannel: wordsChannel,
		options:      options,
	}
}

func (f *FileReader) Prepare(_ context.Context) error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			f.wordsChannel <- line
		}
		file.Close()
		close(f.wordsChannel)
	}()

	return nil
}

func (f *FileReader) Words(_ context.Context) <-chan string {
	out := make(chan string, 1)
	go func() {
		for word := range f.wordsChannel {
			out <- word
		}
		close(out)
	}()
	return out
}
