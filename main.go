package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting the application")
	var textFilePath, url string
	flag.StringVar(&textFilePath, "text", "", "Path to the text file to be processed")
	flag.StringVar(&textFilePath, "t", "", "Path to the text file to be processed")

	flag.StringVar(&url, "url", "", "URL to the text file to be processed")
	flag.StringVar(&url, "u", "", "URL to the text file to be processed")
}
