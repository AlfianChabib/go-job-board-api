package validator

import (
	"net/url"
	"reflect"
	"strings"
	"time"

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
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.RegisterValidation("rfc3339", ValidateRFC3339)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidation("nullable_url", ValidateNullableURL)
	if err != nil {
		panic(err)
	}

	return &structValidator{
		validate:   validate,
		translator: trans,
	}
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

// FormatValidationErrors converts errors into a translated map[string]string
func (v *structValidator) FormatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorsMap[e.Field()] = e.Translate(v.translator)
		}
	}

	return errorsMap
}

// ValidateRFC3339 checks whether a string has a valid ISO 8601 / RFC3339 format.
func ValidateRFC3339(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	// If empty, leave the 'required' and 'required if' tags to throw an error.
	// We only validate the format if the data is actually filled in.
	if dateStr == "" {
		return true
	}

	// Try parsing the string to RFC3339 format
	_, err := time.Parse(time.RFC3339, dateStr)

	// If err is nil (success), return true. If it fails, return false.
	return err == nil
}

func ValidateNullableURL(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	// Jika isinya string kosong "", biarkan lolos
	if val == "" {
		return true
	}

	// Jika ada isinya, pastikan itu adalah URL yang valid
	u, err := url.ParseRequestURI(val)

	// Pastikan tidak error, memiliki skema (http/https), dan ada host-nya
	return err == nil && u.Scheme != "" && u.Host != ""
}
