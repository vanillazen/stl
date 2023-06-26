package errors

import "fmt"

type (
	errors struct{}
)

func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func Wrapf(err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a)
	return Wrap(err, msg)
}
