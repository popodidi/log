package iowriter

import (
	"fmt"
	"io"

	"github.com/ttacon/chalk"

	"github.com/popodidi/log"
)

var (
	styleMap = map[log.Level]chalk.Style{
		log.Debug:    chalk.ResetColor.NewStyle(),
		log.Info:     chalk.Green.NewStyle().WithTextStyle(chalk.Bold),
		log.Notice:   chalk.Cyan.NewStyle().WithTextStyle(chalk.Bold),
		log.Warn:     chalk.Yellow.NewStyle().WithTextStyle(chalk.Bold),
		log.Error:    chalk.Red.NewStyle().WithTextStyle(chalk.Bold),
		log.Critical: chalk.Magenta.NewStyle().WithTextStyle(chalk.Bold),
	}

	timeStyle = chalk.ResetColor.NewStyle().WithTextStyle(chalk.Inverse)
)

const timeFormat = "2006-01-02 15:04:05.000"

// Config defines the writer handler config.
type Config struct {
	Writer     io.Writer
	WithColor  bool
	TimeFormat string
}

// New returns a writer handler with config.
func New(conf Config) log.Handler {
	h := &handler{
		Config: conf,
	}
	if h.TimeFormat == "" {
		h.TimeFormat = timeFormat
	}
	return h
}

type handler struct {
	Config
}

func (h *handler) Handle(entry *log.Entry) {
	tsRaw := entry.Time.Format(h.TimeFormat)
	svRaw := fmt.Sprintf("%s", entry.Level.String())
	tagRaw := fmt.Sprintf("%s", entry.Tag)

	var b []byte
	if !h.WithColor {
		b = []byte(fmt.Sprintf("%s %5s | %s | %s\n", tsRaw, svRaw, tagRaw, entry.Log))
	} else {
		style := styleMap[entry.Level]
		timestamp := timeStyle.Style(tsRaw)
		content := style.Style(fmt.Sprintf("%5s | %s | %s", svRaw, tagRaw, entry.Log))

		b = []byte(fmt.Sprintf("%s %s\n", timestamp, content))
	}
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
