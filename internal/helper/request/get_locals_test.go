package request_test

import (
	"AlfianChabib/go-job-board-api/internal/helper/request"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAuthLocals(t *testing.T) {
	expectedUserID := uuid.New()
	expectedRole := domain.RoleCandidate

	t.Run("Success when userId and role are valid", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			c.Locals("role", expectedRole)

			uid, role, err := request.GetAuthLocals(c)
			require.NoError(t, err)
			assert.Equal(t, expectedUserID, uid)
			assert.Equal(t, expectedRole, role)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Error when userId is missing or nil UUID", func(t *testing.T) {
		app := fiber.New()
		app.Get("/missing", func(c fiber.Ctx) error {
			_, _, err := request.GetAuthLocals(c)
			return err
		})
		app.Get("/nil-uuid", func(c fiber.Ctx) error {
			c.Locals("userId", uuid.Nil)
			c.Locals("role", expectedRole)
			_, _, err := request.GetAuthLocals(c)
			return err
		})

		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/missing", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/nil-uuid", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)
	})

	t.Run("Error when role is missing or empty", func(t *testing.T) {
		app := fiber.New()
		app.Get("/missing-role", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			_, _, err := request.GetAuthLocals(c)
			return err
		})
		app.Get("/empty-role", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			c.Locals("role", domain.UserRole(""))
			_, _, err := request.GetAuthLocals(c)
			return err
		})

		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/missing-role", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/empty-role", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)
	})
}

func TestGetLocalSession(t *testing.T) {
	expectedUserID := uuid.New()
	expectedRole := domain.RoleRecruiter

	t.Run("Success when session locals are valid", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			c.Locals("role", expectedRole)

			session, err := request.GetLocalSession(c)
			require.NoError(t, err)
			assert.Equal(t, expectedUserID, session.UserId)
			assert.Equal(t, expectedRole, session.Role)
			return c.SendStatus(fiber.StatusOK)
		})

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/test", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Error when userId is missing or nil", func(t *testing.T) {
		app := fiber.New()
		app.Get("/no-user", func(c fiber.Ctx) error {
			_, err := request.GetLocalSession(c)
			return err
		})
		app.Get("/nil-user", func(c fiber.Ctx) error {
			c.Locals("userId", uuid.Nil)
			c.Locals("role", expectedRole)
			_, err := request.GetLocalSession(c)
			return err
		})

		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/no-user", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/nil-user", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)
	})

	t.Run("Error when role is missing or empty", func(t *testing.T) {
		app := fiber.New()
		app.Get("/no-role", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			_, err := request.GetLocalSession(c)
			return err
		})
		app.Get("/empty-role", func(c fiber.Ctx) error {
			c.Locals("userId", expectedUserID)
			c.Locals("role", domain.UserRole(""))
			_, err := request.GetLocalSession(c)
			return err
		})

		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/no-role", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/empty-role", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)
	})
}
