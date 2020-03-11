package main

import (
	"time"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/iowriter"
	"github.com/popodidi/log/handlers/iowriter/file"
)

func main() {
	rotateFile := file.Rotate(
		file.PrefixSuffix("log/example-log-", ".txt", file.SecondRotator(1<<7)),
	)
	defer rotateFile.Close()

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler: iowriter.New(iowriter.Config{
			Writer: rotateFile,
			Codec:  iowriter.DefaultCodec(false),
		}),
	})
	logger := log.New("example-log")

	logger.Debug("Debug at %s", time.Now())
	logger.Info("Info at %s", time.Now())
	logger.Notice("Notice at %s", time.Now())
	time.Sleep(time.Second)
	logger.Warn("Warn at %s", time.Now())
	logger.Error("Error at %s", time.Now())
	logger.Critical("Critical at %s", time.Now())
}
