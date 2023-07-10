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

var (
	Empty = Error{
		err:        nil,
		context:    "",
		stackTrace: "",
	}
)

func New(msg string) Error {
	err := errors.New(msg)
	return Wrap(err, "")
}

func Newf(format string, ctxValues ...interface{}) Error {
	msg := fmt.Sprintf(format, ctxValues...)
	err := errors.New(msg)
	return Wrap(err, "")
}

func Wrap(err error, context ...string) Error {
	var ctx string
	if len(context) > 0 {
		ctx = context[0]
	}

	stackTrace := extractStacktrace()
	return Error{
		err:        err,
		context:    ctx,
		stackTrace: stackTrace,
	}
}

func Wrapf(err error, format string, ctxValues ...interface{}) error {
	context := fmt.Sprintf(format, ctxValues...)
	return Wrap(err, context)
}

type Error struct {
	err        error
	context    string
	stackTrace string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %v\n%s", err.context, err.err, err.stackTrace)
}

func (err Error) Unwrap() error {
	return err.err
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
	if wrappedErr, ok := err.(*Error); ok {
		return wrappedErr.Error()
	}
	return err.Error()
}
