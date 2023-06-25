package http

import "github.com/vanillazen/stl/backend/internal/sys/errors"

type (
	HTTPError struct {
		errors.Err
	}
)

func NewError(msg string) HTTPError {
	return HTTPError{
		Err: errors.NewError(msg),
	}
}

var (
	MethodNotAllowedErr   = errors.NewError("method not allowed")
	InvalidResourceErr    = errors.NewError("invalid resource")
	NoUserErr             = errors.NewError("not a valid user in session")
	NoResourceErr         = errors.NewError("no resource ID provided")
	NoAssetReqErr         = errors.NewError("no asset request provided")
	InvalidRequestErr     = errors.NewError("invalid request")
	InvalidRequestDataErr = errors.NewError("invalid request data")
	InvalidJSONBodyErr    = errors.NewError("invalid JSON body")
	InvalidValueTypeErr   = errors.NewError("invalid value type")
	ListNotFoundErr       = errors.NewError("list not found")
)
