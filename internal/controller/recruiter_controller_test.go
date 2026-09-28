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
)

func TestRecruiterController_CompanyCRUD(t *testing.T) {
	recruiterID := uuid.New()
	svc := new(mockRecruiterService)
	ctrl := NewRecruiterController(svc)

	app := setupTestApp()
	app.Post("/unauth", ctrl.CreateCompany)
	app.Get("/unauth", ctrl.GetCompany)
	app.Put("/unauth", ctrl.UpdateCompany)

	app.Post("/company", withSession(recruiterID, domain.RoleRecruiter), ctrl.CreateCompany)
	app.Get("/company", withSession(recruiterID, domain.RoleRecruiter), ctrl.GetCompany)
	app.Put("/company", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateCompany)

	t.Run("Unauthorized checks", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPost, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r2.StatusCode)
		r3, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r3.StatusCode)
	})

	t.Run("CreateCompany validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPost, "/company", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.CreateCompanyRequest{
			Name:         "Acme Corp",
			Location:     "Jakarta",
			Industry:     "IT",
			EmployeeSize: domain.EmployeeSize1To50,
		})
		svc.On("CreateCompany", mock.Anything, recruiterID, mock.AnythingOfType("web.CreateCompanyRequest")).
			Return(nil, errs.ErrCompanyAlreadyExists).Once()
		req2 := httptest.NewRequest(http.MethodPost, "/company", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusConflict, r2.StatusCode)

		svc.On("CreateCompany", mock.Anything, recruiterID, mock.AnythingOfType("web.CreateCompanyRequest")).
			Return(&web.CompanyResponse{ID: uuid.New(), Name: "Acme Corp"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPost, "/company", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusCreated, r3.StatusCode)
	})

	t.Run("GetCompany service error and success", func(t *testing.T) {
		svc.On("GetCompany", mock.Anything, recruiterID).Return(nil, errs.ErrCompanyNotFound).Once()
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/company", nil))
		assert.Equal(t, fiber.StatusNotFound, r1.StatusCode)

		svc.On("GetCompany", mock.Anything, recruiterID).Return(&web.CompanyResponse{Name: "Acme Corp"}, nil).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/company", nil))
		assert.Equal(t, fiber.StatusOK, r2.StatusCode)
	})

	t.Run("UpdateCompany validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPut, "/company", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.UpdateCompanyRequest{
			Name:         "Acme Updated",
			Location:     "Bandung",
			Industry:     "IT",
			EmployeeSize: domain.EmployeeSize51To200,
		})
		svc.On("UpdateCompany", mock.Anything, recruiterID, mock.AnythingOfType("web.UpdateCompanyRequest")).
			Return(nil, errs.ErrCompanyNotFound).Once()
		req2 := httptest.NewRequest(http.MethodPut, "/company", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("UpdateCompany", mock.Anything, recruiterID, mock.AnythingOfType("web.UpdateCompanyRequest")).
			Return(&web.CompanyResponse{Name: "Acme Updated"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPut, "/company", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})
}

func TestRecruiterController_LogoAndBanner(t *testing.T) {
	recruiterID := uuid.New()
	svc := new(mockRecruiterService)
	ctrl := NewRecruiterController(svc)

	app := setupTestApp()
	app.Put("/unauth-logo", ctrl.UpdateLogo)
	app.Put("/unauth-banner", ctrl.UpdateBanner)
	app.Put("/logo", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateLogo)
	app.Put("/banner", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateBanner)

	t.Run("UpdateLogo validation and success branches", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth-logo", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)

		// Missing file
		r2, _ := app.Test(httptest.NewRequest(http.MethodPut, "/logo", nil))
		assert.Equal(t, fiber.StatusBadRequest, r2.StatusCode)

		// Invalid filename
		b3, ct3, _ := createMultipartPayload("file", "bad name.png", "image/png", []byte("img"))
		req3 := httptest.NewRequest(http.MethodPut, "/logo", b3)
		req3.Header.Set("Content-Type", ct3)
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusBadRequest, r3.StatusCode)

		// Invalid extension
		b4, ct4, _ := createMultipartPayload("file", "logo.gif", "image/png", []byte("img"))
		req4 := httptest.NewRequest(http.MethodPut, "/logo", b4)
		req4.Header.Set("Content-Type", ct4)
		r4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusBadRequest, r4.StatusCode)

		// Invalid Content-Type
		b5, ct5, _ := createMultipartPayload("file", "logo.png", "application/pdf", []byte("img"))
		req5 := httptest.NewRequest(http.MethodPut, "/logo", b5)
		req5.Header.Set("Content-Type", ct5)
		r5, _ := app.Test(req5)
		assert.Equal(t, fiber.StatusBadRequest, r5.StatusCode)

		// Service error via fallback "logo" field
		svc.On("UploadLogo", mock.Anything, mock.AnythingOfType("web.UpdateCompanyLogoRequest")).
			Return(nil, errs.ErrUploadLogoFailed).Once()
		b6, ct6, _ := createMultipartPayload("logo", "logo.png", "image/png", []byte("img"))
		req6 := httptest.NewRequest(http.MethodPut, "/logo", b6)
		req6.Header.Set("Content-Type", ct6)
		r6, _ := app.Test(req6)
		assert.Equal(t, fiber.StatusInternalServerError, r6.StatusCode)

		// Success
		svc.On("UploadLogo", mock.Anything, mock.AnythingOfType("web.UpdateCompanyLogoRequest")).
			Return(&web.UploadCompanyLogoResponse{LogoUrl: "http://example.com/logo.png"}, nil).Once()
		b7, ct7, _ := createMultipartPayload("file", "logo.png", "image/png", []byte("img"))
		req7 := httptest.NewRequest(http.MethodPut, "/logo", b7)
		req7.Header.Set("Content-Type", ct7)
		r7, _ := app.Test(req7)
		assert.Equal(t, fiber.StatusOK, r7.StatusCode)
	})

	t.Run("UpdateBanner validation and success branches", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth-banner", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)

		// Missing file
		r2, _ := app.Test(httptest.NewRequest(http.MethodPut, "/banner", nil))
		assert.Equal(t, fiber.StatusBadRequest, r2.StatusCode)

		// Invalid filename
		b3, ct3, _ := createMultipartPayload("file", "bad name.png", "image/png", []byte("img"))
		req3 := httptest.NewRequest(http.MethodPut, "/banner", b3)
		req3.Header.Set("Content-Type", ct3)
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusBadRequest, r3.StatusCode)

		// Invalid extension
		b4, ct4, _ := createMultipartPayload("file", "banner.gif", "image/png", []byte("img"))
		req4 := httptest.NewRequest(http.MethodPut, "/banner", b4)
		req4.Header.Set("Content-Type", ct4)
		r4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusBadRequest, r4.StatusCode)

		// Invalid Content-Type
		b5, ct5, _ := createMultipartPayload("file", "banner.png", "application/pdf", []byte("img"))
		req5 := httptest.NewRequest(http.MethodPut, "/banner", b5)
		req5.Header.Set("Content-Type", ct5)
		r5, _ := app.Test(req5)
		assert.Equal(t, fiber.StatusBadRequest, r5.StatusCode)

		// Service error via fallback "banner" field
		svc.On("UploadBanner", mock.Anything, mock.AnythingOfType("web.UpdateCompanyBannerRequest")).
			Return(nil, errs.ErrUploadBannerFailed).Once()
		b6, ct6, _ := createMultipartPayload("banner", "banner.webp", "image/webp", []byte("img"))
		req6 := httptest.NewRequest(http.MethodPut, "/banner", b6)
		req6.Header.Set("Content-Type", ct6)
		r6, _ := app.Test(req6)
		assert.Equal(t, fiber.StatusInternalServerError, r6.StatusCode)

		// Success
		svc.On("UploadBanner", mock.Anything, mock.AnythingOfType("web.UpdateCompanyBannerRequest")).
			Return(&web.UploadCompanyBannerResponse{BannerUrl: "http://example.com/banner.webp"}, nil).Once()
		b7, ct7, _ := createMultipartPayload("file", "banner.webp", "image/webp", []byte("img"))
		req7 := httptest.NewRequest(http.MethodPut, "/banner", b7)
		req7.Header.Set("Content-Type", ct7)
		r7, _ := app.Test(req7)
		assert.Equal(t, fiber.StatusOK, r7.StatusCode)
	})
}
