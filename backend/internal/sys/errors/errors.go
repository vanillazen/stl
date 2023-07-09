package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	stackSize   = 10
	skipNTraces = 4
)

func New(msg string) error {
	err := errors.New(msg)
	return Wrap(err, "")
}

func Newf(format string, ctxValues ...interface{}) error {
	msg := fmt.Sprintf(format, ctxValues...)
	err := errors.New(msg)
	return Wrap(err, "")
}

func Wrap(err error, context string) error {
	stackTrace := extractStacktrace()
	return &wrappedError{
		err:        err,
		context:    context,
		stackTrace: stackTrace,
	}
}

func Wrapf(err error, format string, ctxValues ...interface{}) error {
	context := fmt.Sprintf(format, ctxValues...)
	return Wrap(err, context)
}

type wrappedError struct {
	err        error
	context    string
	stackTrace string
}

func (we *wrappedError) Error() string {
	return fmt.Sprintf("%s: %v\n%s", we.context, we.err, we.stackTrace)
}

func (we *wrappedError) Unwrap() error {
	return we.err
}

func extractStacktrace() string {
	stack := make([]uintptr, stackSize)
	length := runtime.Callers(skipNTraces, stack)
	stack = stack[:length]
	frames := runtime.CallersFrames(stack)

	var builder strings.Builder
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		fmt.Fprintf(&builder, "\tat %s:%d: %s\n", frame.File, frame.Line, frame.Function)
	}

	return builder.String()
}

func Stacktrace(err error) string {
	if wrappedErr, ok := err.(*wrappedError); ok {
		return wrappedErr.Error()
	}
	return err.Error()
}
