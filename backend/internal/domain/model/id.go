package model

import "github.com/vanillazen/stl/backend/internal/sys/uuid"

type (
	ID struct {
		val uuid.UUID
	}
)

func NewID(uid uuid.UUID) ID {
	return ID{val: uid}
}

func (i *ID) GenID(id ...uuid.UUID) error {
	if i.val != uuid.Nil {
		return nil // already has a value assigned
	}

	if len(id) > 0 {
		i.val = id[0] // If value is provided, use it
		return nil
	}

	val, err := uuid.New()
	if err != nil {
		return err
	}

	i.val = val

	return nil
}

func (i *ID) Val() uuid.UUID {
	return i.val
}

func (i *ID) String() string {
	return i.val.Val
}
