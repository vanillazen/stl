package port

import (
	"context"
)

type (
	// Mailer WIP
	// Bare interface
	Mailer interface {
		SendMail(ctx context.Context, e Email) error
	}

	EmailAddress string

	// Email WIP
	// Not final structure
	Email struct {
		From EmailAddress
		To   []EmailAddress
		CC   []EmailAddress
		BC   []EmailAddress
		Body []byte
	}
)
