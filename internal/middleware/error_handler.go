package middleware

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	customValidator "AlfianChabib/go-job-board-api/pkg/validator"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	var e *fiber.Error
	if errors.As(err, &e) && e != nil {
		code = e.Code
		message = e.Message
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
		return response.ValidationError(c, "Validation failed", out)
	}

	if message != "" && err != nil {
		return response.Error(c, code, message)
	}

	return nil
}

func NewCustomErrorHandler(v customValidator.StructValidator) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		log.Error(err)
		if valErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			formattedErrors := v.FormatValidationErrors(valErrs)
			return response.ValidationError(c, "Validation failed", formattedErrors)
		}

		// 2. Cek apakah error dari HTTP Fiber bawaan (misal 404, 405)
		if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
			return response.Error(c, fiberErr.Code, fiberErr.Message)
		}

		// 3. Fallback untuk internal server error (500)
		return response.Error(c, fiber.StatusInternalServerError, "Internal server error")
	}
}
