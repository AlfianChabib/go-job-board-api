package middleware_test

import (
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockJwtManager struct {
	validateAccessFunc  func(tokenString string) (*domain.JwtCustomClaims, error)
	validateRefreshFunc func(tokenString string) (*domain.JwtCustomClaims, error)
	generateFunc        func(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error)
}

func (m *mockJwtManager) GenerateTokenPair(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error) {
	if m.generateFunc != nil {
		return m.generateFunc(userID, role)
	}
	return nil, nil
}

func (m *mockJwtManager) ValidateAccessToken(tokenString string) (*domain.JwtCustomClaims, error) {
	if m.validateAccessFunc != nil {
		return m.validateAccessFunc(tokenString)
	}
	return nil, nil
}

func (m *mockJwtManager) ValidateRefreshToken(tokenString string) (*domain.JwtCustomClaims, error) {
	if m.validateRefreshFunc != nil {
		return m.validateRefreshFunc(tokenString)
	}
	return nil, nil
}

func TestProtectedMiddleware(t *testing.T) {
	expectedUserID := uuid.New()
	expectedRole := domain.RoleCandidate

	jwtMgr := &mockJwtManager{
		validateAccessFunc: func(tokenString string) (*domain.JwtCustomClaims, error) {
			if tokenString == "valid-token" {
				return &domain.JwtCustomClaims{
					UserID: expectedUserID,
					Role:   expectedRole,
				}, nil
			}
			return nil, errors.New("Invalid token")
		},
	}

	mw := middleware.NewMiddleware(jwtMgr)
	app := fiber.New()
	app.Get("/protected", mw.Protected(), func(c fiber.Ctx) error {
		uid := c.Locals("userId").(uuid.UUID)
		role := c.Locals("role").(domain.UserRole)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user_id": uid.String(),
			"role":    string(role),
		})
	})

	t.Run("Missing Authorization header returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Malformed Authorization header (not 2 parts) returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "BearerOnlyOnePart")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Valid token sets locals and calls Next", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestRequireRolesMiddleware(t *testing.T) {
	mw := middleware.NewMiddleware(&mockJwtManager{})

	app := fiber.New()
	app.Get("/no-role", mw.RequireRoles(domain.RoleRecruiter), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/candidate-trying-recruiter", func(c fiber.Ctx) error {
		c.Locals("role", domain.RoleCandidate)
		return c.Next()
	}, mw.RequireRoles(domain.RoleRecruiter), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/recruiter-allowed", func(c fiber.Ctx) error {
		c.Locals("role", domain.RoleRecruiter)
		return c.Next()
	}, mw.RequireRoles(domain.RoleRecruiter, domain.RoleCandidate), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	t.Run("Missing role in session returns 401", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/no-role", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Disallowed role returns 403 Forbidden", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/candidate-trying-recruiter", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Allowed role returns 200 OK", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/recruiter-allowed", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
