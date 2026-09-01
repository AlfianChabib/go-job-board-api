package validator

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	validate *validator.Validate
}

func NewValidator() *StructValidator {
	return &StructValidator{validate: validator.New()}
}

// Validator needs to implement the Validate method
func (v *StructValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
