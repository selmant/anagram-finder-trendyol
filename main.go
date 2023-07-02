package main

import (
	"context"
	"flag"

	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/app/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting the application")

	cfg := buildConfig()

	application := app.NewAnagramApplication(cfg)
	err := application.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

//nolint:gomnd // magic numbers are used for default values
func buildConfig() config.Config {
	var textFilePath, url, redisHost, redisPassword, storageType string
	var redisPort, redisDB, workerPoolSize, wordsChannelSize int
	flag.StringVar(&textFilePath, "file", "",
		"Path to the text file to be processed. It is required if url is not given")
	flag.StringVar(&url, "url", "",
		"URL to the text file to be processed. It is required if file is not given")
	flag.StringVar(&redisHost, "redis-host", "",
		"Redis host")
	flag.IntVar(&redisPort, "redis-port", 0,
		"Redis port")
	flag.StringVar(&redisPassword, "redis-password", "",
		"Redis password")
	flag.IntVar(&redisDB, "redis-db", 0,
		"Redis db for storage. The given db will be flushed before the application starts")

	// default values
	flag.StringVar(&storageType, "storage-type", "local",
		"Storage type (local or redis). Default is local")
	flag.IntVar(&workerPoolSize, "worker-pool-size", 16,
		"Worker pool size. Default is 16")
	flag.IntVar(&wordsChannelSize, "words-channel-size", 8,
		"Words channel size. Default is 8")

	flag.Parse()

	if textFilePath != "" && url != "" {
		log.Fatal("Please provide only one of the text file path or the URL")
	} else if textFilePath == "" && url == "" {
		log.Fatal("Please provide a text file path or a URL")
	}

	if storageType != "local" && storageType != "redis" {
		log.Fatal("Storage type must be either local or redis")
	}

	cfg := config.Config{
		Input: struct {
			File struct{ Path string }
			URL  struct{ URL string }
		}{
			File: struct{ Path string }{Path: textFilePath},
			URL:  struct{ URL string }{URL: url},
		},
		Redis: struct {
			Host     string
			Port     int
			Password string
			DB       int
		}{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
			DB:       redisDB,
		},
		StorageType:      config.StorageType(storageType),
		WorkerPoolSize:   workerPoolSize,
		WordsChannelSize: wordsChannelSize,
	}
	return cfg
}
