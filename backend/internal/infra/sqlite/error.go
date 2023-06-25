package sqlite

import "github.com/vanillazen/stl/backend/internal/sys/errors"

var (
	NoConnectionError    = errors.NewError("no connection error")
	InvalidResourceIDErr = errors.NewError("invalid resource ID")
	UserNotFoundErr      = errors.NewError("user not found")
)
