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

func TestApplicationController(t *testing.T) {
	candidateID := uuid.New()
	recruiterID := uuid.New()
	jobID := uuid.New()
	appID := uuid.New()

	svc := new(mockApplicationService)
	ctrl := NewApplicationController(svc)

	app := setupTestApp()
	app.Post("/unauth/jobs/:id/applications", ctrl.ApplyJob)
	app.Get("/unauth/jobs/:id/applications", ctrl.GetJobApplications)
	app.Get("/unauth/candidate/applications", ctrl.GetCandidateApplications)
	app.Get("/unauth/applications/:id", ctrl.GetApplicationById)
	app.Patch("/unauth/applications/:id/status", ctrl.UpdateStatus)
	app.Patch("/unauth/candidate/applications/:id/withdraw", ctrl.WithdrawApplication)

	app.Post("/jobs/:id/applications", withSession(candidateID, domain.RoleCandidate), ctrl.ApplyJob)
	app.Get("/jobs/:id/applications", withSession(recruiterID, domain.RoleRecruiter), ctrl.GetJobApplications)
	app.Get("/candidate/applications", withSession(candidateID, domain.RoleCandidate), ctrl.GetCandidateApplications)
	app.Get("/applications/:id", withSession(candidateID, domain.RoleCandidate), ctrl.GetApplicationById)
	app.Patch("/applications/:id/status", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateStatus)
	app.Patch("/candidate/applications/:id/withdraw", withSession(candidateID, domain.RoleCandidate), ctrl.WithdrawApplication)

	t.Run("Unauthorized checks", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPost, "/unauth/jobs/"+jobID.String()+"/applications", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth/jobs/"+jobID.String()+"/applications", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r2.StatusCode)
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth/candidate/applications", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r3.StatusCode)
		r4, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth/applications/"+appID.String(), nil))
		assert.Equal(t, fiber.StatusUnauthorized, r4.StatusCode)
		r5, _ := app.Test(httptest.NewRequest(http.MethodPatch, "/unauth/applications/"+appID.String()+"/status", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r5.StatusCode)
		r6, _ := app.Test(httptest.NewRequest(http.MethodPatch, "/unauth/candidate/applications/"+appID.String()+"/withdraw", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r6.StatusCode)
	})

	t.Run("ApplyJob validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPost, "/jobs/"+jobID.String()+"/applications", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.ApplyJobRequest{
			ResumeUrl: "https://example.com/cv.pdf",
		})
		svc.On("ApplyJob", mock.Anything, candidateID, mock.AnythingOfType("web.ApplyJobRequest")).
			Return(nil, errs.ErrAlreadyApplied).Once()
		req2 := httptest.NewRequest(http.MethodPost, "/jobs/"+jobID.String()+"/applications", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusConflict, r2.StatusCode)

		svc.On("ApplyJob", mock.Anything, candidateID, mock.AnythingOfType("web.ApplyJobRequest")).
			Return(&web.CandidateApplicationResponse{ID: appID, Status: "APPLIED"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPost, "/jobs/"+jobID.String()+"/applications", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusCreated, r3.StatusCode)
	})

	t.Run("GetJobApplications validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.String()+"/applications?status=BAD", nil))
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		svc.On("GetJobApplications", mock.Anything, recruiterID, mock.AnythingOfType("web.GetJobApplicationsRequest")).
			Return(nil, errs.ErrJobForbidden).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.String()+"/applications?status=APPLIED", nil))
		assert.Equal(t, fiber.StatusForbidden, r2.StatusCode)

		svc.On("GetJobApplications", mock.Anything, recruiterID, mock.AnythingOfType("web.GetJobApplicationsRequest")).
			Return(&web.RecruiterApplicationPaginationResponse{Total: 1}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.String()+"/applications?status=APPLIED", nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("GetCandidateApplications validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/candidate/applications?status=BAD", nil))
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		svc.On("GetCandidateApplications", mock.Anything, candidateID, mock.AnythingOfType("web.GetCandidateApplicationsRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/candidate/applications?status=APPLIED", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r2.StatusCode)

		svc.On("GetCandidateApplications", mock.Anything, candidateID, mock.AnythingOfType("web.GetCandidateApplicationsRequest")).
			Return(&web.CandidateApplicationPaginationResponse{Total: 1}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/candidate/applications?status=APPLIED", nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("GetApplicationById invalid UUID, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/applications/not-a-uuid", nil))
		assert.Equal(t, fiber.StatusBadRequest, r1.StatusCode)

		svc.On("GetApplicationById", mock.Anything, candidateID, string(domain.RoleCandidate), appID).
			Return(nil, errs.ErrApplicationNotFound).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/applications/"+appID.String(), nil))
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("GetApplicationById", mock.Anything, candidateID, string(domain.RoleCandidate), appID).
			Return(&web.CandidateApplicationResponse{ID: appID}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/applications/"+appID.String(), nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("UpdateStatus validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPatch, "/applications/"+appID.String()+"/status", bytes.NewBufferString(`{"status":"BAD"}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		svc.On("UpdateStatus", mock.Anything, recruiterID, mock.AnythingOfType("web.UpdateApplicationStatusRequest")).
			Return(nil, errs.ErrApplicationForbidden).Once()
		req2 := httptest.NewRequest(http.MethodPatch, "/applications/"+appID.String()+"/status", bytes.NewBufferString(`{"status":"SHORTLISTED"}`))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusForbidden, r2.StatusCode)

		svc.On("UpdateStatus", mock.Anything, recruiterID, mock.AnythingOfType("web.UpdateApplicationStatusRequest")).
			Return(&web.RecruiterApplicationResponse{ID: appID, Status: "SHORTLISTED"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPatch, "/applications/"+appID.String()+"/status", bytes.NewBufferString(`{"status":"SHORTLISTED"}`))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("WithdrawApplication invalid UUID, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPatch, "/candidate/applications/not-a-uuid/withdraw", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r1.StatusCode)

		svc.On("WithdrawApplication", mock.Anything, candidateID, mock.AnythingOfType("web.WithdrawApplicationRequest")).
			Return(nil, errs.ErrApplicationAlreadyWithdrawn).Once()
		req2 := httptest.NewRequest(http.MethodPatch, "/candidate/applications/"+appID.String()+"/withdraw", bytes.NewBufferString(`{}`))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusBadRequest, r2.StatusCode)

		svc.On("WithdrawApplication", mock.Anything, candidateID, mock.AnythingOfType("web.WithdrawApplicationRequest")).
			Return(&web.CandidateApplicationResponse{ID: appID, Status: "WITHDRAWN"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPatch, "/candidate/applications/"+appID.String()+"/withdraw", bytes.NewBufferString(`{}`))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})
}
