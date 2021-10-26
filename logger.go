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

func (l *logger) WithTag(tags ...string) Logger {
	cloned := l.Clone().(*logger)
	if cloned.conf.Tag == l.id {
		cloned.conf.Tag = strings.Join(tags, ".")
	} else {
		cloned.conf.Tag += "." + strings.Join(tags, ".")
	}
	return cloned
}

func (l *logger) WithLabel(key, value string) Logger {
	cloned := l.Clone().(*logger)
	cloned.GetLabels().Set(key, value)
	return cloned
}

func (l *logger) WithHandler(handlers ...Handler) Logger {
	cloned := l.Clone().(*logger)
	handlers = append(handlers, cloned.conf.Handler)
	cloned.conf.Handler = MultiHandler(handlers...)
	return cloned
}

func (l *logger) GetID() string {
	return l.id
}

func (l *logger) GetTag() string {
	return l.conf.Tag
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
	entry := &Entry{
		Tag:    l.conf.Tag,
		Level:  level,
		Labels: l.labels,
		Log:    fmt.Sprintf(format, args...),
		Time:   time.Now(),
	}
	if level <= Error {
		entry.SetDebugInfo(3)
	}
	l.Handle(entry)
}

func (l *logger) Handle(entry *Entry) {
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

func (l *logger) Log(entry *Entry) {
	l.Handle(entry)
}
