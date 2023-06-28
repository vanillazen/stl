package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
)

type Logger interface {
	SetLogLevel(level LogLevel)
	Debug(v ...any)
	Debugf(format string, a ...any)
	Info(v ...any)
	Infof(format string, a ...any)
	Error(v ...any)
	Errorf(format string, a ...any)
}

type SimpleLogger struct {
	debug    *log.Logger
	info     *log.Logger
	error    *log.Logger
	logLevel LogLevel
}

func NewLogger(logLevel string) *SimpleLogger {
	level := ToValidLevel(logLevel)
	return &SimpleLogger{
		debug:    log.New(os.Stdout, "[DBG] ", log.LstdFlags),
		info:     log.New(os.Stdout, "[INF] ", log.LstdFlags),
		error:    log.New(os.Stderr, "[ERR] ", log.LstdFlags),
		logLevel: level,
	}
}

func (l *SimpleLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *SimpleLogger) Debug(v ...any) {
	if l.logLevel <= Debug {
		l.debug.Println(v...)
	}
}

func (l *SimpleLogger) Debugf(format string, a ...any) {
	if l.logLevel <= Debug {
		message := fmt.Sprintf(format, a...)
		l.debug.Println(message)
	}
}

func (l *SimpleLogger) Info(v ...any) {
	if l.logLevel <= Info {
		l.info.Println(v...)
	}
}

func (l *SimpleLogger) Infof(format string, a ...any) {
	if l.logLevel <= Info {
		message := fmt.Sprintf(format, a...)
		l.info.Println(message)
	}
}

func (l *SimpleLogger) Error(v ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprint(v...)
		l.error.Println(message)
	}
}

func (l *SimpleLogger) Errorf(format string, a ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprintf(format, a...)
		l.error.Println(message)
	}
}

func ToValidLevel(level string) LogLevel {
	level = strings.ToLower(level)

	switch level {
	case "debug", "dbg":
		return Debug
	case "info", "inf":
		return Info
	case "error", "err":
		return Error
	default:
		return Error
	}
}

// SetDebugOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetDebugOutput(debug *bytes.Buffer) {
	sl.debug = log.New(debug, "", 0)
}

// SetInfoOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetInfoOutput(info *bytes.Buffer) {
	sl.info = log.New(info, "", 0)
}

// SetErrorOutput set the internal logger.
// Used for package testing.
func (sl *SimpleLogger) SetErrorOutput(error *bytes.Buffer) {
	sl.error = log.New(error, "", 0)
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
