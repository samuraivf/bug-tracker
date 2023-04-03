package handler

//go:generate mockgen -source=validator.go -destination=mocks/validator.go
type Validator interface {
	Struct(s interface{}) error
}

type CustomValidator struct {
	validator Validator
}

func newValidator(validator Validator) *CustomValidator {
	return &CustomValidator{validator}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
