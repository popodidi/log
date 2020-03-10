package main

import (
	"os"
	"time"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/iowriter"
	"github.com/popodidi/log/handlers/iowriter/file"
)

func main() {
	singleFile, err := file.Single("log")
	if err != nil {
		os.Exit(1)
	}
	defer singleFile.Close()

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler: iowriter.New(iowriter.Config{
			Writer:    singleFile,
			WithColor: false,
		}),
	})
	logger := log.New("example-log")

	logger.Debug("Debug at %s", time.Now())
	logger.Info("Info at %s", time.Now())
	logger.Notice("Notice at %s", time.Now())
	logger.Warn("Warn at %s", time.Now())
	logger.Error("Error at %s", time.Now())
	logger.Critical("Critical at %s", time.Now())
}
