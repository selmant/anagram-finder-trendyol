package app

import (
	"context"
	"fmt"

	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
)

type AnagramApplication struct {
	AnagramInput   input.DataReader
	AnagramStorage storage.AnagramStorage
}

func NewAnagramApplication(input input.DataReader, storage storage.AnagramStorage) *AnagramApplication {
	return &AnagramApplication{
		AnagramInput:   input,
		AnagramStorage: storage,
	}
}

func (app *AnagramApplication) Run() error {
	err := app.AnagramInput.Prepare(context.Background())
	if err != nil {
		return err
	}

	for word := range app.AnagramInput.Lines(context.Background()) {
		fmt.Println(word)
		/*err = app.AnagramStorage.Store(word)
		if err != nil {
			return err
		}*/
	}

	return nil
}
