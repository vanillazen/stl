package log

import (
	"bytes"
	"fmt"
)

type TestLogger struct {
	debug    *bytes.Buffer
	info     *bytes.Buffer
	error    *bytes.Buffer
	logLevel LogLevel
}

func NewTestLogger(logLevel string) *TestLogger {
	level := ToValidLevel(logLevel)
	return &TestLogger{
		debug:    bytes.NewBuffer(nil),
		info:     bytes.NewBuffer(nil),
		error:    bytes.NewBuffer(nil),
		logLevel: level,
	}
}

func (l *TestLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *TestLogger) Debug(v ...interface{}) {
	if l.logLevel <= Debug {
		message := fmt.Sprintln(v...)
		l.debug.WriteString(message)
	}
}

func (l *TestLogger) Debugf(format string, a ...interface{}) {
	if l.logLevel <= Debug {
		message := fmt.Sprintf(format, a...)
		l.debug.WriteString(message)
	}
}

func (l *TestLogger) Info(v ...interface{}) {
	if l.logLevel <= Info {
		message := fmt.Sprintln(v...)
		l.info.WriteString(message)
	}
}

func (l *TestLogger) Infof(format string, a ...interface{}) {
	if l.logLevel <= Info {
		message := fmt.Sprintf(format, a...)
		l.info.WriteString(message)
	}
}

func (l *TestLogger) Error(v ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprintln(v...)
		l.error.WriteString(message)
	}
}

func (l *TestLogger) Errorf(format string, a ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprintf(format, a...)
		l.error.WriteString(message)
	}
}

// GetDebugLogs returns the captured debug logs.
func (l *TestLogger) GetDebugLogs() string {
	return l.debug.String()
}

// GetInfoLogs returns the captured info logs.
func (l *TestLogger) GetInfoLogs() string {
	return l.info.String()
}

// GetErrorLogs returns the captured error logs.
func (l *TestLogger) GetErrorLogs() string {
	return l.error.String()
}
