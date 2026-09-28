package validator_test

import (
	"AlfianChabib/go-job-board-api/pkg/validator"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleValidationStruct struct {
	Name      string `json:"name,omitempty" validate:"required,min=3"`
	StartDate string `json:"start_date" validate:"rfc3339"`
	Website   string `json:"website" validate:"nullable_url"`
	Ignored   string `json:"-" validate:"required"`
}

func TestStructValidator(t *testing.T) {
	v := validator.NewValidator()
	require.NotNil(t, v)

	t.Run("Valid struct with empty optional custom fields passes", func(t *testing.T) {
		s := sampleValidationStruct{
			Name:      "Alice",
			StartDate: "",
			Website:   "",
			Ignored:   "present",
		}
		err := v.Validate(s)
		assert.NoError(t, err)
	})

	t.Run("Valid struct with populated RFC3339 and nullable_url passes", func(t *testing.T) {
		s := sampleValidationStruct{
			Name:      "Alice",
			StartDate: "2026-09-28T10:00:00Z",
			Website:   "https://example.com/path",
			Ignored:   "present",
		}
		err := v.Validate(s)
		assert.NoError(t, err)
	})

	t.Run("Invalid fields fail validation and format properly", func(t *testing.T) {
		s := sampleValidationStruct{
			Name:      "Al",
			StartDate: "not-a-date",
			Website:   "invalid-url",
			Ignored:   "",
		}
		err := v.Validate(s)
		require.Error(t, err)

		formatted := v.FormatValidationErrors(err)
		assert.Contains(t, formatted, "name")
		assert.Contains(t, formatted, "start_date")
		assert.Contains(t, formatted, "website")
		// Ignored has json:"-" so RegisterTagNameFunc returns "", which falls back to struct field name "Ignored"
		assert.Contains(t, formatted, "Ignored")
	})

	t.Run("ValidateNullableURL rejects URL without scheme or host", func(t *testing.T) {
		s := sampleValidationStruct{
			Name:    "Alice",
			Website: "http://",
			Ignored: "ok",
		}
		err := v.Validate(s)
		require.Error(t, err)
		formatted := v.FormatValidationErrors(err)
		assert.Contains(t, formatted, "website")
	})

	t.Run("FormatValidationErrors with non-validation error returns empty map", func(t *testing.T) {
		formatted := v.FormatValidationErrors(errors.New("some standard error"))
		assert.Empty(t, formatted)
	})
}
