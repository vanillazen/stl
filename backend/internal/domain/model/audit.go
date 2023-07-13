package model

import "time"

type (
	Audit struct {
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func NewAudit(createdAt, updatedAt time.Time) Audit {
	return Audit{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
