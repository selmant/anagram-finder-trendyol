package main

import (
	"context"
	"flag"

	"github.com/redis/go-redis/v9"
	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting the application")
	var textFilePath, url string
	flag.StringVar(&textFilePath, "text", "", "Path to the text file to be processed")
	flag.StringVar(&textFilePath, "t", "", "Path to the text file to be processed")

	flag.StringVar(&url, "url", "", "URL to the text file to be processed")
	flag.StringVar(&url, "u", "", "URL to the text file to be processed")

	flag.Parse()

	fileReader := input.NewFileReader("test-anagrams.txt", input.DefaultFileReaderOptions())
	localStorage := storage.NewLocalStorage()

	redisClient := redis.NewClient(&redis.Options{})
	redisStorage := storage.NewRedisStorage(*redisClient)
	_ = app.NewAnagramApplication(&fileReader, localStorage)
	redisApp := app.NewAnagramApplication(&fileReader, redisStorage)
	err := redisApp.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
