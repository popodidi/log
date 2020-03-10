package main

import (
	"os"
	"time"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/iowriter"
	"github.com/popodidi/log/handlers/iowriter/file"
	"github.com/popodidi/log/handlers/multi"
)

func main() {
	singleFile, err := file.Single("log")
	if err != nil {
		os.Exit(1)
	}
	defer singleFile.Close()
	fileHandler := iowriter.New(iowriter.Config{
		Writer:    singleFile,
		WithColor: false,
	})
	stdOutHandler := iowriter.New(iowriter.Config{
		Writer:    os.Stdout,
		WithColor: true,
	})

	// multi handler
	handler := multi.New(fileHandler, stdOutHandler)

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler:   handler,
	})
	logger := log.New("example-log")

	logger.Debug("Debug at %s", time.Now())
	logger.Info("Info at %s", time.Now())
	logger.Notice("Notice at %s", time.Now())
	logger.Warn("Warn at %s", time.Now())
	logger.Error("Error at %s", time.Now())
	logger.Critical("Critical at %s", time.Now())
}
