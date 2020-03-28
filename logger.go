package log

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var defaultConf = Config{
	Threshold: Debug,
	Handler:   (*null)(nil),
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

const (
	idKey = "github.com/popodidi/log.id"
)

// Config defines a config for Logger.
type Config struct {
	Tag       string
	Threshold Level
	Handler   Handler
}

var _ Logger = (*logger)(nil)

type logger struct {
	conf   Config
	id     string
	labels Labels
}

// Set sets the default config of logger.
func Set(conf Config) {
	defaultConf = conf
}

// New returns a new Logger with default config.
func New(tag ...string) Logger {
	conf := defaultConf
	conf.Tag = strings.Join(tag, ".")
	return NewLogger(conf)
}

// NewLogger returns a new Logger with config.
func NewLogger(config Config) Logger {
	logger := &logger{
		conf:   config,
		id:     getID(),
		labels: NewLabels(),
	}
	if len(logger.conf.Tag) == 0 {
		logger.conf.Tag = logger.id
	}
	logger.labels.Set(idKey, logger.id)
	return logger
}

func getID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).
			Intn(len(letterRunes))]
	}
	return string(b)
}

func (l *logger) Clone() Logger {
	cloned := &logger{
		conf:   l.conf,
		id:     l.id,
		labels: l.labels.Clone(),
	}
	return cloned
}

func (l *logger) GetID() string {
	return l.id
}

func (l *logger) GetLabels() Labels {
	return l.labels
}

// nolint: goprintffuncname
func (l *logger) Debug(format string, args ...interface{}) {
	l.handle(Debug, format, args...)
}

// nolint: goprintffuncname
func (l *logger) Info(format string, args ...interface{}) {
	l.handle(Info, format, args...)
}

// nolint: goprintffuncname
func (l *logger) Notice(format string, args ...interface{}) {
	l.handle(Notice, format, args...)
}

// nolint: goprintffuncname
func (l *logger) Warn(format string, args ...interface{}) {
	l.handle(Warn, format, args...)
}

// nolint: goprintffuncname
func (l *logger) Error(format string, args ...interface{}) {
	l.handle(Error, format, args...)
}

// nolint: goprintffuncname
func (l *logger) Critical(format string, args ...interface{}) {
	l.handle(Critical, format, args...)
}

// nolint: goprintffuncname
func (l *logger) handle(level Level, format string, args ...interface{}) {
	l.Log(&Entry{
		StackNum: 6,
		Tag:      l.conf.Tag,
		Level:    level,
		Labels:   l.labels,
		Log:      fmt.Sprintf(format, args...),
		Time:     time.Now(),
	})
}

func (l *logger) Log(entry *Entry) {
	if entry.Level > l.conf.Threshold {
		return
	}

	l.conf.Handler.Handle(entry)
	if entry.Level <= Critical {
		if closer, ok := l.conf.Handler.(CloseHandler); ok {
			err := closer.Close()
			if err != nil {
				log.Println("failed to close handler. err:", err)
			}
		}
		os.Exit(1)
	}
}
