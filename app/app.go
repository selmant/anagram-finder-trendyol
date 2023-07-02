package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	storagelib "github.com/selmant/anagram-finder-trendyol/internal/storage"
	log "github.com/sirupsen/logrus"
)

type AnagramApplication struct {
	Input          input.DataReader
	AnagramStorage storagelib.Storage
}

func NewAnagramApplication(cfg config.Config) *AnagramApplication {
	log.Info("Redis storage is selected")
	config.GlobalConfig = &cfg
	var reader input.DataReader
	if cfg.Input.File.Path != "" {
		reader = input.NewFileReader(cfg.Input.File.Path)
	} else {
		reader = input.NewURLReader(cfg.Input.URL.URL)
	}

	var storage storagelib.Storage

	if cfg.StorageType == config.StorageTypeLocal {
		storage = storagelib.NewLocalStorage()
	} else {
		log.Info("Redis storage is selected")
		redisClient := storagelib.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
		if cmd := redisClient.Ping(context.Background()); cmd.Err() != nil {
			log.Fatal(cmd.Err())
		}
		log.Info("Redis client created")
		storage = storagelib.NewRedisStorage(redisClient)
		log.Info("Redis storage created")
	}

	return &AnagramApplication{
		Input:          reader,
		AnagramStorage: storage,
	}
}

func (app *AnagramApplication) Run(ctx context.Context) error {
	err := app.Input.Prepare(context.Background())
	if err != nil {
		return err
	}

	err = app.hashAndStore(ctx)
	if err != nil {
		// TODO: handle error
		log.Error(err)
	}
	err = app.printAnagrams(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (app *AnagramApplication) hashAndStore(ctx context.Context) error {
	start := time.Now()
	job := internal.NewReadAndMatchAnagramJob(app.AnagramStorage, app.Input)
	pool := internal.NewWorkerPool(config.GlobalConfig.WorkerPoolSize, job)
	err := pool.Start(ctx)
	if err != nil {
		return err
	}
	log.Info("Hash and store took ", time.Since(start))
	return nil
}

func (app *AnagramApplication) printAnagrams(ctx context.Context) error {
	errSlice := make([]error, 0)
	start := time.Now()
	all, errs := app.AnagramStorage.AllAnagrams(ctx)
	for all != nil || errs != nil {
		select {
		case anagrams, ok := <-all:
			if !ok {
				all = nil
			} else if len(anagrams) > 1 {
				//nolint:forbidigo // It's ok
				fmt.Println(strings.Join(anagrams, ", "))
			}
		case chanErr, ok := <-errs:
			if !ok {
				errs = nil
			} else {
				errSlice = append(errSlice, chanErr)
			}
		}
	}
	if len(errSlice) > 0 {
		return errors.Join(errSlice...)
	}
	log.Info("Read and print took ", time.Since(start))
	return nil
}
