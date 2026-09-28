package response_test

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseHelpers(t *testing.T) {
	app := fiber.New()

	type samplePayload struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	app.Get("/success", func(c fiber.Ctx) error {
		return response.Success(c, fiber.StatusCreated, "Created resource", samplePayload{ID: 1, Name: "Job"})
	})

	app.Get("/ok", func(c fiber.Ctx) error {
		return response.OK(c, "OK message", samplePayload{ID: 2, Name: "Candidate"})
	})

	app.Get("/message", func(c fiber.Ctx) error {
		return response.Message(c, fiber.StatusOK, "Simple message")
	})

	app.Get("/error", func(c fiber.Ctx) error {
		return response.Error(c, fiber.StatusBadRequest, "Bad request error")
	})

	app.Get("/validation-error", func(c fiber.Ctx) error {
		details := map[string]string{"email": "email is required"}
		return response.ValidationError(c, "Validation failed", details)
	})

	t.Run("Success returns custom status and payload", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/success", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var parsed response.Response[samplePayload]
		require.NoError(t, json.Unmarshal(body, &parsed))
		assert.True(t, parsed.Success)
		assert.Equal(t, "Created resource", parsed.Message)
		require.NotNil(t, parsed.Data)
		assert.Equal(t, 1, parsed.Data.ID)
		assert.Equal(t, "Job", parsed.Data.Name)
	})

	t.Run("OK returns 200 status and payload", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/ok", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var parsed response.Response[samplePayload]
		require.NoError(t, json.Unmarshal(body, &parsed))
		assert.True(t, parsed.Success)
		assert.Equal(t, "OK message", parsed.Message)
		require.NotNil(t, parsed.Data)
		assert.Equal(t, 2, parsed.Data.ID)
	})

	t.Run("Message returns SimpleResponse", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/message", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var parsed response.SimpleResponse
		require.NoError(t, json.Unmarshal(body, &parsed))
		assert.True(t, parsed.Success)
		assert.Equal(t, "Simple message", parsed.Message)
	})

	t.Run("Error returns ErrorResponse with success=false", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var parsed response.ErrorResponse[any]
		require.NoError(t, json.Unmarshal(body, &parsed))
		assert.False(t, parsed.Success)
		assert.Equal(t, "Bad request error", parsed.Message)
	})

	t.Run("ValidationError returns 422 with details", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/validation-error", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var parsed response.ErrorResponse[map[string]string]
		require.NoError(t, json.Unmarshal(body, &parsed))
		assert.False(t, parsed.Success)
		assert.Equal(t, "Validation failed", parsed.Message)
		assert.Equal(t, "email is required", parsed.Errors["email"])
	})
}
