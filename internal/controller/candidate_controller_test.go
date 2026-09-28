package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCandidateController_ProfileAndSkills(t *testing.T) {
	userID := uuid.New()

	t.Run("Get unauthorized, service error, and success", func(t *testing.T) {
		svc := new(mockCandidateService)
		ctrl := NewCandidateController(svc)
		app := setupTestApp()
		app.Get("/unauth", ctrl.Get)
		app.Get("/candidate", withSession(userID, domain.RoleCandidate), ctrl.Get)

		// 1. Unauthorized
		resp1, err := app.Test(httptest.NewRequest(http.MethodGet, "/unauth", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		// 2. Service error
		svc.On("Get", mock.Anything, userID).Return(nil, errs.ErrNotFound).Once()
		resp2, err := app.Test(httptest.NewRequest(http.MethodGet, "/candidate", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp2.StatusCode)

		// 3. Success
		svc.On("Get", mock.Anything, userID).Return(&web.GetCandidateResponse{UserId: userID}, nil).Once()
		resp3, err := app.Test(httptest.NewRequest(http.MethodGet, "/candidate", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})

	t.Run("Update unauthorized, bind error, service error, and success", func(t *testing.T) {
		svc := new(mockCandidateService)
		ctrl := NewCandidateController(svc)
		app := setupTestApp()
		app.Put("/unauth", ctrl.Update)
		app.Put("/candidate", withSession(userID, domain.RoleCandidate), ctrl.Update)

		// 1. Unauthorized
		resp1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		// 2. Validation error
		req2 := httptest.NewRequest(http.MethodPut, "/candidate", bytes.NewBufferString(`{}`))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp2.StatusCode)

		// 3. Service error
		payload, _ := json.Marshal(web.UpdateCandidateRequest{
			Headline: "Go Engineer",
			Phone:    "+628123456789",
		})
		svc.On("Update", mock.Anything, mock.AnythingOfType("domain.Profile")).
			Return(nil, errs.ErrInternalServer).Once()
		req3 := httptest.NewRequest(http.MethodPut, "/candidate", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		resp3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusInternalServerError, resp3.StatusCode)

		// 4. Success
		svc.On("Update", mock.Anything, mock.AnythingOfType("domain.Profile")).
			Return(&web.UpdateCandidateResponse{UserId: userID}, nil).Once()
		req4 := httptest.NewRequest(http.MethodPut, "/candidate", bytes.NewBuffer(payload))
		req4.Header.Set("Content-Type", "application/json")
		resp4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusOK, resp4.StatusCode)
	})

	t.Run("UpdateSkills unauthorized, bind error, service error, and success", func(t *testing.T) {
		svc := new(mockCandidateService)
		ctrl := NewCandidateController(svc)
		app := setupTestApp()
		app.Put("/unauth", ctrl.UpdateSkills)
		app.Put("/skills", withSession(userID, domain.RoleCandidate), ctrl.UpdateSkills)

		resp1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodPut, "/skills", bytes.NewBufferString(`{}`))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp2.StatusCode)

		payload, _ := json.Marshal(web.UpdateCandidateSkillsRequest{
			Skills: []web.SkillRequest{{Name: "Go"}},
		})
		svc.On("UpdateSkills", mock.Anything, userID, mock.AnythingOfType("web.UpdateCandidateSkillsRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		req3 := httptest.NewRequest(http.MethodPut, "/skills", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		resp3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusInternalServerError, resp3.StatusCode)

		svc.On("UpdateSkills", mock.Anything, userID, mock.AnythingOfType("web.UpdateCandidateSkillsRequest")).
			Return([]domain.Skill{{Name: "go"}}, nil).Once()
		req4 := httptest.NewRequest(http.MethodPut, "/skills", bytes.NewBuffer(payload))
		req4.Header.Set("Content-Type", "application/json")
		resp4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusOK, resp4.StatusCode)
	})
}

func TestCandidateController_Avatar(t *testing.T) {
	userID := uuid.New()

	t.Run("UpdateAvatar all validation branches and success", func(t *testing.T) {
		svc := new(mockCandidateService)
		ctrl := NewCandidateController(svc)
		app := setupTestApp()
		app.Put("/unauth", ctrl.UpdateAvatar)
		app.Put("/avatar", withSession(userID, domain.RoleCandidate), ctrl.UpdateAvatar)

		// 1. Unauthorized
		resp1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		// 2. Missing file
		resp2, _ := app.Test(httptest.NewRequest(http.MethodPut, "/avatar", nil))
		assert.Equal(t, fiber.StatusBadRequest, resp2.StatusCode)

		// 3. Invalid filename (contains space)
		body3, ct3, _ := createMultipartPayload("avatar", "bad name.png", "image/png", []byte("img"))
		req3 := httptest.NewRequest(http.MethodPut, "/avatar", body3)
		req3.Header.Set("Content-Type", ct3)
		resp3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusBadRequest, resp3.StatusCode)

		// 4. Invalid extension (.gif)
		body4, ct4, _ := createMultipartPayload("avatar", "avatar.gif", "image/png", []byte("img"))
		req4 := httptest.NewRequest(http.MethodPut, "/avatar", body4)
		req4.Header.Set("Content-Type", ct4)
		resp4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusBadRequest, resp4.StatusCode)

		// 5. Invalid Content-Type
		body5, ct5, _ := createMultipartPayload("avatar", "avatar.png", "application/pdf", []byte("img"))
		req5 := httptest.NewRequest(http.MethodPut, "/avatar", body5)
		req5.Header.Set("Content-Type", ct5)
		resp5, _ := app.Test(req5)
		assert.Equal(t, fiber.StatusBadRequest, resp5.StatusCode)

		// 6. Service error
		svc.On("UploadAvatar", mock.Anything, mock.AnythingOfType("web.UpdateCandidateAvatarRequest")).
			Return(nil, errors.New("upload err")).Once()
		body6, ct6, _ := createMultipartPayload("avatar", "avatar.png", "image/png", []byte("img"))
		req6 := httptest.NewRequest(http.MethodPut, "/avatar", body6)
		req6.Header.Set("Content-Type", ct6)
		resp6, _ := app.Test(req6)
		assert.Equal(t, fiber.StatusInternalServerError, resp6.StatusCode)

		// 7. Success
		url := "http://localhost:9000/avatar/avatar.png"
		svc.On("UploadAvatar", mock.Anything, mock.AnythingOfType("web.UpdateCandidateAvatarRequest")).
			Return(&url, nil).Once()
		body7, ct7, _ := createMultipartPayload("avatar", "avatar.png", "image/png", []byte("img"))
		req7 := httptest.NewRequest(http.MethodPut, "/avatar", body7)
		req7.Header.Set("Content-Type", ct7)
		resp7, _ := app.Test(req7)
		assert.Equal(t, fiber.StatusOK, resp7.StatusCode)
	})

	t.Run("DeleteAvatar unauthorized, service error, and success", func(t *testing.T) {
		svc := new(mockCandidateService)
		ctrl := NewCandidateController(svc)
		app := setupTestApp()
		app.Delete("/unauth", ctrl.DeleteAvatar)
		app.Delete("/avatar", withSession(userID, domain.RoleCandidate), ctrl.DeleteAvatar)

		resp1, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/unauth", nil))
		assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

		svc.On("DeleteAvatar", mock.Anything, userID).Return(errs.ErrInternalServer).Once()
		resp2, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/avatar", nil))
		assert.Equal(t, fiber.StatusInternalServerError, resp2.StatusCode)

		svc.On("DeleteAvatar", mock.Anything, userID).Return(nil).Once()
		resp3, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/avatar", nil))
		assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
	})
}

func TestCandidateController_Experiences(t *testing.T) {
	userID := uuid.New()
	expID := uuid.New()
	endDate := "2026-01-01T00:00:00Z"

	svc := new(mockCandidateService)
	ctrl := NewCandidateController(svc)
	app := setupTestApp()
	app.Get("/unauth-exp", ctrl.GetExperiences)
	app.Post("/unauth-exp", ctrl.CreateExperience)
	app.Put("/unauth-exp/:experienceId", ctrl.UpdateExperience)
	app.Delete("/unauth-exp/:experienceId", ctrl.DeleteExperience)

	app.Get("/exp", withSession(userID, domain.RoleCandidate), ctrl.GetExperiences)
	app.Post("/exp", withSession(userID, domain.RoleCandidate), ctrl.CreateExperience)
	app.Put("/exp/:experienceId", withSession(userID, domain.RoleCandidate), ctrl.UpdateExperience)
	app.Delete("/exp/:experienceId", withSession(userID, domain.RoleCandidate), ctrl.DeleteExperience)

	t.Run("Unauthorized checks", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth-exp", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)
		r2, _ := app.Test(httptest.NewRequest(http.MethodPost, "/unauth-exp", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r2.StatusCode)
		r3, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth-exp/"+expID.String(), nil))
		assert.Equal(t, fiber.StatusUnauthorized, r3.StatusCode)
		r4, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/unauth-exp/"+expID.String(), nil))
		assert.Equal(t, fiber.StatusUnauthorized, r4.StatusCode)
	})

	t.Run("GetExperiences error and success", func(t *testing.T) {
		svc.On("GetExperiences", mock.Anything, userID).Return(nil, errs.ErrInternalServer).Once()
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/exp", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r1.StatusCode)

		svc.On("GetExperiences", mock.Anything, userID).Return([]domain.Experience{}, nil).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/exp", nil))
		assert.Equal(t, fiber.StatusOK, r2.StatusCode)
	})

	t.Run("CreateExperience validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPost, "/exp", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.CreateExperienceRequest{
			CompanyName: "Acme",
			Position:    "Dev",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     &endDate,
		})
		svc.On("CreateExperience", mock.Anything, userID, mock.AnythingOfType("web.CreateExperienceRequest")).
			Return(errs.ErrInternalServer).Once()
		req2 := httptest.NewRequest(http.MethodPost, "/exp", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusInternalServerError, r2.StatusCode)

		svc.On("CreateExperience", mock.Anything, userID, mock.AnythingOfType("web.CreateExperienceRequest")).
			Return(nil).Once()
		req3 := httptest.NewRequest(http.MethodPost, "/exp", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("UpdateExperience validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPut, "/exp/"+expID.String(), bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.UpdateExperienceRequest{
			CompanyName: "Acme",
			Position:    "Lead",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     &endDate,
		})
		svc.On("UpdateExperience", mock.Anything, userID, mock.AnythingOfType("web.UpdateExperienceRequest")).
			Return(errs.ErrExperienceNotFound).Once()
		req2 := httptest.NewRequest(http.MethodPut, "/exp/"+expID.String(), bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("UpdateExperience", mock.Anything, userID, mock.AnythingOfType("web.UpdateExperienceRequest")).
			Return(nil).Once()
		req3 := httptest.NewRequest(http.MethodPut, "/exp/"+expID.String(), bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("DeleteExperience invalid UUID, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/exp/not-a-uuid", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r1.StatusCode)

		svc.On("DeleteExperience", mock.Anything, userID, expID).Return(errs.ErrExperienceNotFound).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/exp/"+expID.String(), nil))
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("DeleteExperience", mock.Anything, userID, expID).Return(nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/exp/"+expID.String(), nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})
}
