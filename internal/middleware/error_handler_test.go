package middleware_test

import (
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/pkg/errs"
	customValidator "AlfianChabib/go-job-board-api/pkg/validator"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type dummyReq struct {
	Email string `json:"email" validate:"required,email"`
}

func TestErrorHandler(t *testing.T) {
	v := customValidator.NewValidator()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/fiber-error", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "bad fiber request")
	})

	app.Get("/validation-error", func(c fiber.Ctx) error {
		return v.Validate(dummyReq{Email: "not-an-email"})
	})

	app.Get("/generic-error", func(c fiber.Ctx) error {
		return errors.New("something crashed")
	})

	app.Get("/nil-error", func(c fiber.Ctx) error {
		return middleware.ErrorHandler(c, nil)
	})

	t.Run("Handles fiber.Error", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/fiber-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Handles validator.ValidationErrors", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/validation-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Handles generic error with 500", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/generic-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Returns nil when err is nil", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/nil-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestNewCustomErrorHandler(t *testing.T) {
	v := customValidator.NewValidator()
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.NewCustomErrorHandler(v),
	})

	app.Get("/app-error", func(c fiber.Ctx) error {
		return errs.ErrJobNotFound
	})

	app.Get("/fiber-error", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTeapot, "I am a teapot")
	})

	app.Get("/gorm-not-found", func(c fiber.Ctx) error {
		return gorm.ErrRecordNotFound
	})

	app.Get("/validation-error", func(c fiber.Ctx) error {
		return v.Validate(dummyReq{Email: ""})
	})

	app.Get("/json-unmarshal-error", func(c fiber.Ctx) error {
		return &json.UnmarshalTypeError{
			Value: "string",
			Type:  reflect.TypeOf(0),
			Field: "min_salary",
		}
	})

	app.Get("/unknown-error", func(c fiber.Ctx) error {
		return errors.New("unexpected db failure")
	})

	t.Run("Handles AppError", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/app-error", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Handles fiber.Error", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/fiber-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusTeapot, resp.StatusCode)
	})

	t.Run("Handles gorm.ErrRecordNotFound", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/gorm-not-found", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Handles validator.ValidationErrors", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/validation-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Handles json.UnmarshalTypeError", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/json-unmarshal-error", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Handles unknown error with 500", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/unknown-error", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
