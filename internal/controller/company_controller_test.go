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

func TestCompanyController(t *testing.T) {
	t.Run("GetCompanies validation error, service error, and success", func(t *testing.T) {
		svc := new(mockCompanyService)
		ctrl := NewCompanyController(svc)
		app := setupTestApp()
		app.Get("/companies", ctrl.GetCompanies)

		// 1. Validation error (limit > 100)
		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies?limit=500", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp1.StatusCode)

		// 2. Service error
		svc.On("GetCompanies", mock.Anything, mock.AnythingOfType("web.GetCompaniesRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies?page=1&limit=10", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp2.StatusCode)

		// 3. Success
		svc.On("GetCompanies", mock.Anything, mock.AnythingOfType("web.GetCompaniesRequest")).
			Return(&web.CompanyPaginationResponse{Total: 1}, nil).Once()
		resp3, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies?page=1&limit=10", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})

	t.Run("GetCompanyById invalid UUID, fallback param companyId, service error, and success", func(t *testing.T) {
		svc := new(mockCompanyService)
		ctrl := NewCompanyController(svc)
		app := setupTestApp()
		app.Get("/companies/:id", ctrl.GetCompanyById)

		// 1. Invalid UUID in URI bind returns 500 from custom error handler
		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies/not-a-uuid", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp1.StatusCode)

		// 2. Service error
		compID := uuid.New()
		svc.On("GetCompanyById", mock.Anything, compID).Return(nil, errs.ErrCompanyNotFound).Once()
		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies/"+compID.String(), nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp2.StatusCode)

		// 3. Success via :id
		svc.On("GetCompanyById", mock.Anything, compID).Return(&web.CompanyDetailResponse{
			CompanyResponse: web.CompanyResponse{ID: compID, Name: "Acme"},
		}, nil).Once()
		resp3, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies/"+compID.String(), nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)

		// 4. Fallback param companyId when StructValidator is not rejecting empty ID on URI bind
		appNoVal := fiber.New(fiber.Config{
			ErrorHandler: setupTestApp().ErrorHandler,
		})
		appNoVal.Get("/alt/:companyId", ctrl.GetCompanyById)
		resp4, err := appNoVal.Test(httptest.NewRequest(http.MethodGet, "/alt/invalid-uuid", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp4.StatusCode)

		svc.On("GetCompanyById", mock.Anything, compID).Return(&web.CompanyDetailResponse{
			CompanyResponse: web.CompanyResponse{ID: compID, Name: "Acme"},
		}, nil).Once()
		resp5, err := appNoVal.Test(httptest.NewRequest(http.MethodGet, "/alt/"+compID.String(), nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp5.StatusCode)
	})
}
