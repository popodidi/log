package stackdriver

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers"
)

// Config defines the stackdriver config.
type Config struct {
	Ctx    context.Context
	Codec  handlers.Codec
	Parent string
	Opts   []option.ClientOption
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
		client: client,
		codec:  conf.Codec,
	}
	if h.codec == nil {
		h.codec = &codec{}
	}
	return h, nil
}

type handler struct {
	client *logging.Client
	codec  handlers.Codec
}

func (h *handler) Close() error {
	return h.client.Close()
}

func (h *handler) Handle(entry *log.Entry) {
	logger := h.client.
		Logger(entry.Tag, logging.CommonLabels(entry.Labels.Map())).
		StandardLogger(severityMap[entry.Level])
	_, err := logger.Writer().Write(h.codec.Encode(entry))
	if err != nil {
		fmt.Println("Failed to write log to stackdriver. err:", err)
	}
}
