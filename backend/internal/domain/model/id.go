package model

import (
	"github.com/vanillazen/stl/backend/internal/sys/uuid"
)

type (
	ID struct {
		UUID uuid.UUID
	}
)

func NewID(uid uuid.UUID) ID {
	return ID{UUID: uid}
}

func (i *ID) GenID(id ...uuid.UUID) error {
	if i.UUID != uuid.Nil {
		return nil // already has a value assigned
	}

	if len(id) > 0 {
		i.UUID = id[0] // If value is provided, use it
		return nil
	}

	i.UUID = uuid.NewUUID()

	return nil
}

func (i *ID) Val() uuid.UUID {
	return i.UUID
}

func (i *ID) String() string {
	return i.UUID.Val
}
