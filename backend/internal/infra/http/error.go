package http

import "github.com/vanillazen/stl/backend/internal/sys/errors"

var (
	MethodNotAllowedErr   = errors.New("method not allowed")
	InvalidResourceErr    = errors.New("invalid resource")
	NoUserErr             = errors.New("not a valid user in session")
	NoResourceErr         = errors.New("no resource ID provided")
	NoAssetReqErr         = errors.New("no asset request provided")
	InvalidRequestErr     = errors.New("invalid request")
	InvalidRequestDataErr = errors.New("invalid request data")
	InvalidJSONBodyErr    = errors.New("invalid JSON body")
	InvalidValueTypeErr   = errors.New("invalid value type")
	ListNotFoundErr       = errors.New("list not found")
)
