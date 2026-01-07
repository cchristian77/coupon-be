package util

import (
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/fhttp"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()
}

func Validate(input any) error {
	if v == nil {
		v = validator.New()
	}

	err := v.Struct(input)

	var fieldErrs []fhttp.OptionalData
	if errors.As(err, &validator.ValidationErrors{}) {
		for _, e := range err.(validator.ValidationErrors) {
			fieldErrs = append(fieldErrs, fhttp.OptionalData{
				Key:   e.Field(),
				Value: e.Error(),
			})
		}
	}

	if err != nil {
		return fhttp.NewErrorResponse(
			http.StatusBadRequest, sharedErrs.ErrKindValidation.String(), "Validation Error", fieldErrs...)
	}

	return nil
}
