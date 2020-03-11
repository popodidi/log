package iowriter

import (
	"fmt"
	"io"
	"os"

	"github.com/popodidi/log"
	"github.com/popodidi/log/handlers"
)

// Config defines the writer handler config.
type Config struct {
	Codec  handlers.Codec
	Writer io.Writer
}

// Stdout returns a handler that encodes with default codec and writes to os.Stdout
func Stdout(color bool) log.Handler {
	return &handler{
		Config: Config{
			Codec:  handlers.DefaultCodec(color),
			Writer: os.Stdout,
		},
	}
}

// New returns a writer handler with config.
func New(conf Config) log.Handler {
	h := &handler{
		Config: conf,
	}
	if h.Writer == nil {
		fmt.Println("no writer found. use os.Stdout")
		h.Writer = os.Stdout
	}
	if h.Codec == nil {
		h.Writer.Write([]byte("not codec found. use log.DefaultCodec"))
		h.Codec = handlers.DefaultCodec(false)
	}
	return h
}

var _ log.CloseHandler = (*handler)(nil)

type handler struct {
	Config
}

func (h *handler) Handle(entry *log.Entry) {
	b := h.Codec.Encode(entry)
	_, err := h.Writer.Write(b)
	if err != nil {
		fmt.Println("Failed to write log to writer. err:", err)
	}
}

func (h *handler) Close() error {
	if closer, ok := h.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
