package sqlite

import "github.com/vanillazen/stl/backend/internal/sys/errors"

var (
	NoConnectionError    = errors.New("no connection error")
	InvalidResourceIDErr = errors.New("invalid resource ID")
	UserNotFoundErr      = errors.New("user not found")
)
