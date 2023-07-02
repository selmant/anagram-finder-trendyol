package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	log "github.com/sirupsen/logrus"
)

type AnagramApplication struct {
	AnagramInput   input.DataReader
	AnagramStorage storage.Storage
}

func NewAnagramApplication(input input.DataReader, storage storage.Storage) *AnagramApplication {
	return &AnagramApplication{
		AnagramInput:   input,
		AnagramStorage: storage,
	}
}

func (app *AnagramApplication) Run(ctx context.Context) error {
	err := app.AnagramInput.Prepare(context.Background())
	if err != nil {
		return err
	}

	start := time.Now()
	job := internal.NewHashAndStoreAnagramJob(app.AnagramStorage, app.AnagramInput)
	pool := internal.NewWorkerPool(8, job)
	err = pool.Start(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Error(time.Now().Sub(start))
	time.Sleep(time.Second * 3)

	start = time.Now()
	all, errs := app.AnagramStorage.AllAnagrams(ctx)
	for all != nil || errs != nil {
		select {
		case anagrams, ok := <-all:
			if !ok {
				all = nil
			} else {
				if len(anagrams) > 1 {
					fmt.Println(strings.Join(anagrams, ", "))
				}
			}
		case chanErr, ok := <-errs:
			if !ok {
				errs = nil
			} else {
				return chanErr
			}
		}
	}
	log.Error(time.Now().Sub(start))

	return nil
}
