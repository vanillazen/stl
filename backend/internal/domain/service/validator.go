package service

import (
	"errors"

	"github.com/vanillazen/stl/backend/internal/domain/model"
	"github.com/vanillazen/stl/backend/internal/sys/validator"
)

type (
	ListValidator struct {
		validator.Validator
		Model model.List
	}
)

func NewListValidator(m model.List) ListValidator {
	return ListValidator{
		Validator: validator.NewValidator(),
		Model:     m,
	}
}

func (v ListValidator) ValidateForCreate() error {
	// Username
	ok0 := v.ValidateRequiredName()
	ok1 := v.ValidateMinLengthName(2)

	if ok0 && ok1 {
		return nil
	}

	return errors.New("list has errors")
}

func (v ListValidator) ValidateForUpdate() error {
	return errors.New("not implemented yet")
}

func (v ListValidator) ValidateRequiredName(errMsg ...string) (ok bool) {
	list := v.Model

	ok = v.ValidateRequired(list.Name)
	if ok {
		return true
	}

	msg := validator.ValidatorMsg.RequiredErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	v.Errors["Name"] = append(v.Errors["Name"], msg)
	return false
}

func (v ListValidator) ValidateMinLengthName(min int, errMsg ...string) (ok bool) {
	m := v.Model

	ok = v.ValidateMinLength(m.Name, min)
	if ok {
		return true
	}

	msg := validator.ValidatorMsg.MinLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	v.Errors["Name"] = append(v.Errors["Name"], msg)
	return false
}

func (v ListValidator) ValidateMaxLengthName(max int, errMsg ...string) (ok bool) {
	m := v.Model

	ok = v.ValidateMaxLength(m.Name, max)
	if ok {
		return true
	}

	msg := validator.ValidatorMsg.MaxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	v.Errors["Name"] = append(v.Errors["Name"], msg)
	return false
}
