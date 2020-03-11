package iowriter

import (
	"fmt"

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

// DefaultCodec returns the default codec.
func DefaultCodec(withColor bool) Codec {
	c := &defaultCodec{
		WithColor:  withColor,
		TimeFormat: timeFormat,
	}
	return c
}

type defaultCodec struct {
	WithColor  bool
	TimeFormat string
}

func (t *defaultCodec) Encode(entry *log.Entry) []byte {
	tsRaw := entry.Time.Format(t.TimeFormat)
	svRaw := fmt.Sprintf("%s", entry.Level.String())
	tagRaw := fmt.Sprintf("%s", entry.Tag)

	if !t.WithColor {
		return []byte(fmt.Sprintf("%s %5s | %s | %s\n", tsRaw, svRaw, tagRaw, entry.Log))
	}

	style := styleMap[entry.Level]
	timestamp := timeStyle.Style(tsRaw)
	content := style.Style(fmt.Sprintf("%5s | %s | %s", svRaw, tagRaw, entry.Log))
	return []byte(fmt.Sprintf("%s %s\n", timestamp, content))
}
