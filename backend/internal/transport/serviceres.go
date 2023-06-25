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
		msg            string        // Human readable message exposed to client
		exposeInternal bool          // Expose internal error to te client flag
		valErrSet      v.ValErrorSet // Properties validation errors
		err            error         // Internal error
	}
)

func (sr *ServiceRes) Msg() string {
	return sr.msg
}

func (sr *ServiceRes) ValidationErrors() v.ValErrorSet {
	return sr.valErrSet
}

func (sr *ServiceRes) Err() error {
	if sr.exposeInternal {
		return sr.err
	}

	return nil
}

func NewServiceRes(valErrSet v.ValErrorSet, err error, cfg *config.Config) ServiceRes {
	return ServiceRes{
		msg:            validationError,
		exposeInternal: cfg.GetBool(config.Key.APIErrorExposeInt),
		valErrSet:      valErrSet,
		err:            err,
	}
}
