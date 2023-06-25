package log

import (
	"bytes"
	"fmt"
)

type TestLogger struct {
	debugLogs *bytes.Buffer
	infoLogs  *bytes.Buffer
	errorLogs *bytes.Buffer
	logLevel  LogLevel
}

func NewTestLogger(logLevel string) *TestLogger {
	level := ToValidLevel(logLevel)
	return &TestLogger{
		debugLogs: bytes.NewBuffer(nil),
		infoLogs:  bytes.NewBuffer(nil),
		errorLogs: bytes.NewBuffer(nil),
		logLevel:  level,
	}
}

func (l *TestLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *TestLogger) Debug(v ...interface{}) {
	if l.logLevel <= Debug {
		message := fmt.Sprintln(v...)
		l.debugLogs.WriteString(message)
	}
}

func (l *TestLogger) Debugf(format string, a ...interface{}) {
	if l.logLevel <= Debug {
		message := fmt.Sprintf(format, a...)
		l.debugLogs.WriteString(message)
	}
}

func (l *TestLogger) Info(v ...interface{}) {
	if l.logLevel <= Info {
		message := fmt.Sprintln(v...)
		l.infoLogs.WriteString(message)
	}
}

func (l *TestLogger) Infof(format string, a ...interface{}) {
	if l.logLevel <= Info {
		message := fmt.Sprintf(format, a...)
		l.infoLogs.WriteString(message)
	}
}

func (l *TestLogger) Error(v ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprintln(v...)
		l.errorLogs.WriteString(message)
	}
}

func (l *TestLogger) Errorf(format string, a ...interface{}) {
	if l.logLevel <= Error {
		message := fmt.Sprintf(format, a...)
		l.errorLogs.WriteString(message)
	}
}

// GetDebugLogs returns the captured debug logs.
func (l *TestLogger) GetDebugLogs() string {
	return l.debugLogs.String()
}

// GetInfoLogs returns the captured info logs.
func (l *TestLogger) GetInfoLogs() string {
	return l.infoLogs.String()
}

// GetErrorLogs returns the captured error logs.
func (l *TestLogger) GetErrorLogs() string {
	return l.errorLogs.String()
}
