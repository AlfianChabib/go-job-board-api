package middleware

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/pkg/errs"
	customValidator "AlfianChabib/go-job-board-api/pkg/validator"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
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
		code := fiber.StatusInternalServerError
		var message string

		if appErr, ok := errors.AsType[*errs.AppError](err); ok {
			code = appErr.Code
			message = appErr.Message

		} else if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
			code = fiberErr.Code
			message = fiberErr.Message

		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			message = "Data not found"

		} else if valErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			formattedErrors := v.FormatValidationErrors(valErrs)
			return response.ValidationError(c, "Validation failed", formattedErrors)

		} else if jsonErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			code = http.StatusBadRequest
			message = "Data type does not match the field: " + jsonErr.Field

		} else {
			log.Println(err)
			message = "Internal server error"

		}

		return response.Error(c, code, message)
	}
}
