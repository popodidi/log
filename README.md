# Log

`log` is a logging library that supports multiple handlers.

## Usage

```go
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
```

Note that `logger.Critical` calls `os.Exit(1)` after logging.

## Handlers

For now, `log` supports the following handlers

- io.Writer
  - single file
  - rotating files

Checkout `example/` for more details.
