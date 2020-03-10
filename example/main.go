package main

import (
	"os"
	"time"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/iowriter"
)

func main() {
	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler: iowriter.New(iowriter.Config{
			Writer:    os.Stdout,
			WithColor: true,
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
