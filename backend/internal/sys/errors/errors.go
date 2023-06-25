package errors

import "fmt"

type (
	errors struct{}
)

func Wrap(message string, err error) error {
	return fmt.Errorf("%s: %w", message, err)
}

