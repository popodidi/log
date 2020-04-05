package log_test

import (
	"os"
	"time"

	"google.golang.org/api/option"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers/async"
	"github.com/popodidi/log/handlers/codec"
	"github.com/popodidi/log/handlers/filtered"
	"github.com/popodidi/log/handlers/iowriter"
	"github.com/popodidi/log/handlers/iowriter/file"
	"github.com/popodidi/log/handlers/multi"
	"github.com/popodidi/log/handlers/stackdriver"
)

func Example_stdout() {
	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler:   iowriter.Stdout(true),
	})
	logger := log.New("example-log")

	logger.Debug("Debug at %s", time.Now())
	logger.Info("Info at %s", time.Now())
	logger.Notice("Notice at %s", time.Now())
	logger.Warn("Warn at %s", time.Now())
	logger.Error("Error at %s", time.Now())
	logger.Critical("Critical at %s", time.Now())
}

func Example_singleFile() {
	singleFile, err := file.Single("log")
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = singleFile.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler: iowriter.New(iowriter.Config{
			Writer: singleFile,
			Codec:  codec.Default(false),
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

func Example_rotateFile() {
	rotateFile := file.Rotate(
		file.PrefixSuffix("log/example-log-", ".txt", file.SecondRotator(1<<7)),
	)
	defer func() {
		err := rotateFile.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler: iowriter.New(iowriter.Config{
			Writer: rotateFile,
			Codec:  codec.Default(false),
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

func Example_asyncRotateFile() {
	handler := iowriter.New(iowriter.Config{
		Writer: file.Rotate(
			file.PrefixSuffix("log/example-log-", ".txt", file.SecondRotator(1<<7)),
		),
		Codec: codec.Default(false),
	})
	handler = async.New(handler)
	defer func() {
		err := handler.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler:   handler,
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

func Example_multi() {
	singleFile, err := file.Single("log")
	if err != nil {
		os.Exit(1)
	}

	// multi handler
	handler := iowriter.New(iowriter.Config{
		Writer: singleFile,
		Codec:  codec.Default(false),
	})
	handler = multi.New(handler, iowriter.Stdout(true))
	handler = multi.New(handler, iowriter.Stdout(false))
	defer func() {
		err = handler.(log.CloseHandler).Close()
		if err != nil {
			os.Exit(1)
		}
	}()

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

func Example_filtered() {
	// Configure logger
	log.Set(log.Config{
		Threshold: log.Debug,
		Handler:   filtered.Warn(iowriter.Stdout(true)),
	})
	logger := log.New("example-log")

	logger.Debug("Debug at %s", time.Now())
	logger.Info("Info at %s", time.Now())
	logger.Notice("Notice at %s", time.Now())
	logger.Warn("Warn at %s", time.Now())
	logger.Error("Error at %s", time.Now())
	logger.Critical("Critical at %s", time.Now())
}

func Example_stackdriver() {
	handler, err := stackdriver.New(stackdriver.Config{
		LogName: "xxx",
		Parent:  "gcp-project-id",
		Opts: []option.ClientOption{
			option.WithCredentialsFile("path/to/credentials.json"),
		},
	})
	if err != nil {
		os.Exit(1)
	}

	// multi handler
	handler = multi.New(handler, iowriter.Stdout(true))
	defer func() {
		err = handler.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

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
