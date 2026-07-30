package validator

import (
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
)

type Checker struct {
	*validator.Validate
}

func New() *Checker {
	return &Checker{validator.New()}
}

func (v *Checker) Struct(data interface{}) error {
	var result error
	if errs := v.Validate.Struct(data); errs != nil {
		validationErrors, ok := errs.(validator.ValidationErrors)
		if !ok {
			return errors.WithStack(errs)
		}
		for _, err := range validationErrors {
			result = multierror.Append(result, err)
		}
	}
	if result == nil {
		return nil
	}
	return errors.WithStack(result)
}
