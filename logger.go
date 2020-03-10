package log

import (
	"fmt"
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

// Config defines a config for Logger.
type Config struct {
	Tag       string
	Threshold Level
	Handler   Handler
}

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
	id := getID()
	conf := defaultConf
	conf.Tag = strings.Join(tag, ".")
	if len(conf.Tag) == 0 {
		conf.Tag = id
	}
	return &logger{
		conf:   conf,
		id:     id,
		labels: NewLabels(),
	}
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

func (l *logger) Debug(format string, args ...interface{}) {
	l.handle(Debug, format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.handle(Info, format, args...)
}

func (l *logger) Notice(format string, args ...interface{}) {
	l.handle(Notice, format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.handle(Warn, format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.handle(Error, format, args...)
}

func (l *logger) Critical(format string, args ...interface{}) {
	l.handle(Critical, format, args...)
}

func (l *logger) handle(level Level, format string, args ...interface{}) {
	l.Log(&Entry{
		StackNum: 4,
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
	if entry.Level <= Error {
		entry.Log = fmt.Sprintf("%s: %s", entry.DebugInfo(), entry.Log)
	}
	l.conf.Handler.Handle(entry)
	if entry.Level <= Critical {
		os.Exit(1)
	}
}
