package middleware

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) && e != nil {
		code = e.Code
	}

	var validationErrors validator.ValidationErrors
	matchedValidationError := errors.As(err, &validationErrors)
	if matchedValidationError {
		out := make([]fiber.Map, 0, len(validationErrors))
		for _, e := range validationErrors {
			out = append(out, fiber.Map{
				"field": strings.ToLower(e.Field()),
				"rule":  e.Tag(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation Error",
			"errors":  out,
		})
	}

	if e.Message != "" && err != nil {
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"message": e.Message,
		})
	}

	return nil
}
