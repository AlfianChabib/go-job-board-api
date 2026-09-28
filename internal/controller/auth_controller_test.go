package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthController(t *testing.T) {
	t.Run("Register bind error, service error, and success", func(t *testing.T) {
		svc := new(mockAuthService)
		ctrl := NewAuthController(svc, "production")
		app := setupTestApp()
		app.Post("/register", ctrl.Register)

		// 1. Validation / bind error
		req1 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":"invalid"}`))
		req1.Header.Set("Content-Type", "application/json")
		resp1, err := app.Test(req1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		validPayload, _ := json.Marshal(web.RegisterRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     domain.RoleCandidate,
		})
		svc.On("Register", mock.Anything, mock.AnythingOfType("web.RegisterRequest")).
			Return(nil, errs.ErrUserAlreadyExists).Once()

		req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validPayload))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err := app.Test(req2)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, resp2.StatusCode)

		// 3. Success
		svc.On("Register", mock.Anything, mock.AnythingOfType("web.RegisterRequest")).
			Return(&web.RegisterResponse{
				ID:    uuid.New(),
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  domain.RoleCandidate,
			}, nil).Once()

		req3 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validPayload))
		req3.Header.Set("Content-Type", "application/json")
		resp3, err := app.Test(req3)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})

	t.Run("Login bind error, service error, and success with cookie", func(t *testing.T) {
		svc := new(mockAuthService)
		ctrl := NewAuthController(svc, "development")
		app := setupTestApp()
		app.Post("/login", ctrl.Login)

		// 1. Validation error
		req1 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		resp1, err := app.Test(req1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		validPayload, _ := json.Marshal(web.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		})
		svc.On("Login", mock.Anything, mock.AnythingOfType("web.LoginRequest")).
			Return(nil, errs.ErrInvalidCredentials).Once()

		req2 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(validPayload))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err := app.Test(req2)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)

		// 3. Success
		svc.On("Login", mock.Anything, mock.AnythingOfType("web.LoginRequest")).
			Return(&domain.TokenPair{
				AccessToken:  "access-tok",
				RefreshToken: "refresh-tok",
			}, nil).Once()

		req3 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(validPayload))
		req3.Header.Set("Content-Type", "application/json")
		resp3, err := app.Test(req3)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
		assert.NotEmpty(t, resp3.Header.Get("Set-Cookie"))
	})

	t.Run("LogOut missing cookie, service error, and success", func(t *testing.T) {
		svc := new(mockAuthService)
		ctrl := NewAuthController(svc, "development")
		app := setupTestApp()
		app.Post("/logout", ctrl.LogOut)

		// 1. Missing cookie fails validation
		req1 := httptest.NewRequest(http.MethodPost, "/logout", nil)
		resp1, err := app.Test(req1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		svc.On("Logout", mock.Anything, "my-refresh-token").Return(errs.ErrUnauthorized).Once()
		req2 := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req2.AddCookie(&http.Cookie{Name: "refresh_token", Value: "my-refresh-token"})
		resp2, err := app.Test(req2)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)

		// 3. Success
		svc.On("Logout", mock.Anything, "my-refresh-token").Return(nil).Once()
		req3 := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req3.AddCookie(&http.Cookie{Name: "refresh_token", Value: "my-refresh-token"})
		resp3, err := app.Test(req3)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})

	t.Run("RefreshToken bind error, service error, and success", func(t *testing.T) {
		svc := new(mockAuthService)
		ctrl := NewAuthController(svc, "development")
		app := setupTestApp()
		app.Post("/refresh", ctrl.RefreshToken)

		// 1. Missing fields fails validation
		req1 := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		resp1, err := app.Test(req1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		payload, _ := json.Marshal(web.RefreshTokenRequest{
			RefreshToken: "ref-tok",
			AccessToken:  "acc-tok",
		})
		svc.On("RefreshToken", mock.Anything, mock.AnythingOfType("web.RefreshTokenRequest")).
			Return(nil, errs.ErrTokenRevoked).Once()

		req2 := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err := app.Test(req2)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)

		// 3. Success
		svc.On("RefreshToken", mock.Anything, mock.AnythingOfType("web.RefreshTokenRequest")).
			Return(&web.RefreshTokenResponse{
				TokenPair: &domain.TokenPair{
					AccessToken:  "new-acc",
					RefreshToken: "new-ref",
				},
			}, nil).Once()

		req3 := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		resp3, err := app.Test(req3)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})
}
