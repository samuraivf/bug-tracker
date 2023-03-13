package handler

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	validator *validator.Validate
}

func newValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
