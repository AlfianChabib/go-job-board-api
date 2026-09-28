package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataController(t *testing.T) {
	t.Run("GetSkills validation error, service error, and success", func(t *testing.T) {
		svc := new(mockDataService)
		ctrl := NewDataController(svc)
		app := setupTestApp()
		app.Get("/skills", ctrl.GetSkills)

		// 1. Validation error (limit > 100)
		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/skills?limit=200", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		svc.On("GetSkills", mock.Anything, mock.AnythingOfType("web.GetDataRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/skills?search=go&limit=5", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp2.StatusCode)

		// 3. Success
		svc.On("GetSkills", mock.Anything, mock.AnythingOfType("web.GetDataRequest")).
			Return([]web.SkillResponse{{ID: uuid.New(), Name: "go", Label: "Go"}}, nil).Once()
		resp3, err := app.Test(httptest.NewRequest(http.MethodGet, "/skills?search=go&limit=5", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})

	t.Run("GetCurrencyCodes validation error, service error, and success", func(t *testing.T) {
		svc := new(mockDataService)
		ctrl := NewDataController(svc)
		app := setupTestApp()
		app.Get("/currencies", ctrl.GetCurrencyCodes)

		// 1. Validation error (limit > 100)
		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/currencies?limit=200", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		svc.On("GetCurrencyCodes", mock.Anything, mock.AnythingOfType("web.GetDataRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/currencies?search=IDR", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp2.StatusCode)

		// 3. Success
		svc.On("GetCurrencyCodes", mock.Anything, mock.AnythingOfType("web.GetDataRequest")).
			Return([]web.CurrencyResponse{{AlphabeticCode: "IDR", Currency: "Rupiah"}}, nil).Once()
		resp3, err := app.Test(httptest.NewRequest(http.MethodGet, "/currencies?search=IDR", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})
}
