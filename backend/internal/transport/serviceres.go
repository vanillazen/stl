package transport

import (
	"github.com/vanillazen/stl/backend/internal/sys/config"
	v "github.com/vanillazen/stl/backend/internal/sys/validator"
)

const (
	validationError = "Check fields with errors"
)

type (
	ServiceRes struct {
		msg       string        // Human-readable message exposed to client
		valErrSet v.ValErrorSet // Properties validation errors
		err       error         // Internal error
	}
)

func (sr *ServiceRes) Msg() string {
	return sr.msg
}

func (sr *ServiceRes) ValidationErrors() v.ValErrorSet {
	return sr.valErrSet
}

func (sr *ServiceRes) Err() error {
	return sr.err
}

func NewServiceRes(valErrSet v.ValErrorSet, err error, cfg *config.Config) ServiceRes {
	return ServiceRes{
		msg:       validationError,
		valErrSet: valErrSet,
		err:       err,
	}
}
