package stackdriver

import (
	"cloud.google.com/go/logging"

	"github.com/popodidi/log"
)

var severityMap = map[log.Level]logging.Severity{
	log.Critical: logging.Critical,
	log.Error:    logging.Error,
	log.Warn:     logging.Warning,
	log.Notice:   logging.Notice,
	log.Info:     logging.Info,
	log.Debug:    logging.Debug,
}
