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

	// multi handler
	handler := iowriter.New(iowriter.Config{
		Writer: singleFile,
		Codec:  iowriter.DefaultCodec(false),
	})
	handler = multi.New(handler, iowriter.Stdout(true))
	handler = multi.New(handler, iowriter.Stdout(false))

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
