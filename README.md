# Log

![Go Test](https://github.com/popodidi/log/workflows/Go%20Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/popodidi/log)](https://goreportcard.com/report/github.com/popodidi/log)
[![Documentation](https://godoc.org/github.com/popodidi/log?status.svg)](http://godoc.org/github.com/popodidi/log)

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
```

Note that `logger.Critical` calls `os.Exit(1)` after logging.
Checkout the [doc](https://godoc.org/github.com/popodidi/log) for usage examples.

## Handlers

`Handler` interface defines how the log.Entry will be handled.
For the real implementations, we could separate them into are two kinds.

- Intrinsic handlers - handlers that really handle logs
- Wrappers - handlers that simply wraps another handler

For now, `log` supports the following handlers

### Intrinsic handlers

- io.Writer
     - stdout
     - single file
     - rotating files
     - custom `io.Writer`s
- stackdriver

### Wrapper

- filtered
     - level filtering handler
     - custom filter
- multi
     - sync multi handler

## Caveats

For now, the `multi` handler runs the underlying handlers synchronously. The
file writers of `handler/iowriter/file` also runs synchronously. Use these
handlers with concerns in performance-critical code.
