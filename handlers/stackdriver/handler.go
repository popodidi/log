package stackdriver

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers"
	"github.com/popodidi/log/handlers/codec"
)

const tagKey = "github.com/popodidi/log/handlers/stackdriver.tag"

// Config defines the stackdriver config.
type Config struct {
	Ctx     context.Context
	Codec   handlers.Codec
	Parent  string
	Opts    []option.ClientOption
	LogName string
}

// New returns a new stackdriver handler.
func New(conf Config) (log.CloseHandler, error) {
	ctx := conf.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := logging.NewClient(ctx, conf.Parent, conf.Opts...)
	if err != nil {
		return nil, err
	}
	h := &handler{
		logName: conf.LogName,
		client:  client,
		codec:   conf.Codec,
	}
	if h.codec == nil {
		h.codec = codec.Simple()
	}
	return h, nil
}

type handler struct {
	logName string
	client  *logging.Client
	codec   handlers.Codec
}

func (h *handler) Close() error {
	return h.client.Close()
}

func (h *handler) Handle(entry *log.Entry) {
	m := entry.Labels.Map()
	m[tagKey] = entry.Tag
	logger := h.client.
		Logger(h.logName, logging.CommonLabels(m)).
		StandardLogger(severityMap[entry.Level])
	_, err := logger.Writer().Write(h.codec.Encode(entry))
	if err != nil {
		fmt.Println("Failed to write log to stackdriver. err:", err)
	}
}
