package middleware

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		out := make([]fiber.Map, 0, len(validationErrors))
		for _, e := range validationErrors {
			// e.Field() - field name, e.Tag() - failed rule,
			// e.Param() - rule parameter, e.Value() - invalid value
			out = append(out, fiber.Map{
				"field": strings.ToLower(e.Field()),
				"rule":  e.Tag(),
				"value": e.Error(),
				"param": e.Param(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": out,
		})
	}
	return nil
}
