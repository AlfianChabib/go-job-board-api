package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type StructValidator interface {
	Validate(out any) error
	FormatValidationErrors(err error) map[string]string
}

type structValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewValidator() StructValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english)
	trans, _ := uni.GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &structValidator{
		validate:   validate,
		translator: trans,
	}
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

// FormatValidationErrors mengubah error menjadi map[string]string yang sudah diterjemahkan
func (v *structValidator) FormatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorsMap[e.Field()] = e.Translate(v.translator)
		}
	}

	return errorsMap
}
