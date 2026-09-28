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

func TestJobController(t *testing.T) {
	recruiterID := uuid.New()
	jobID := uuid.New()
	svc := new(mockJobService)
	ctrl := NewJobController(svc)

	app := setupTestApp()
	app.Post("/unauth-jobs", ctrl.CreateJob)
	app.Get("/unauth-recruiter-jobs", ctrl.GetRecruiterJobs)
	app.Put("/unauth-jobs/:id", ctrl.UpdateJob)
	app.Patch("/unauth-jobs/:id/status", ctrl.UpdateJobStatus)
	app.Delete("/unauth-jobs/:id", ctrl.DeleteJob)

	app.Get("/jobs", ctrl.GetJobs)
	app.Get("/jobs/:id", ctrl.GetJobById)
	app.Post("/jobs", withSession(recruiterID, domain.RoleRecruiter), ctrl.CreateJob)
	app.Get("/recruiter/jobs", withSession(recruiterID, domain.RoleRecruiter), ctrl.GetRecruiterJobs)
	app.Put("/jobs/:id", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateJob)
	app.Patch("/jobs/:id/status", withSession(recruiterID, domain.RoleRecruiter), ctrl.UpdateJobStatus)
	app.Delete("/jobs/:id", withSession(recruiterID, domain.RoleRecruiter), ctrl.DeleteJob)

	t.Run("Unauthorized checks", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPost, "/unauth-jobs", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r1.StatusCode)
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/unauth-recruiter-jobs", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r2.StatusCode)
		r3, _ := app.Test(httptest.NewRequest(http.MethodPut, "/unauth-jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusUnauthorized, r3.StatusCode)
		r4, _ := app.Test(httptest.NewRequest(http.MethodPatch, "/unauth-jobs/"+jobID.String()+"/status", nil))
		assert.Equal(t, fiber.StatusUnauthorized, r4.StatusCode)
		r5, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/unauth-jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusUnauthorized, r5.StatusCode)
	})

	t.Run("CreateJob validation error, service error, and success", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBufferString(`{}`))
		req1.Header.Set("Content-Type", "application/json")
		r1, _ := app.Test(req1)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		payload, _ := json.Marshal(web.CreateJobRequest{
			Title:          "Backend Developer",
			Description:    "Build scalable Go services",
			EmploymentType: "FULL_TIME",
			WorkMode:       "REMOTE",
			Location:       "Jakarta",
		})
		svc.On("CreateJob", mock.Anything, recruiterID, mock.AnythingOfType("web.CreateJobRequest")).
			Return(nil, errs.ErrCompanyNotFound).Once()
		req2 := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(payload))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("CreateJob", mock.Anything, recruiterID, mock.AnythingOfType("web.CreateJobRequest")).
			Return(&web.JobResponse{ID: jobID, Title: "Backend Developer"}, nil).Once()
		req3 := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusCreated, r3.StatusCode)
	})

	t.Run("GetJobs validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs?limit=500", nil))
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		svc.On("GetJobs", mock.Anything, mock.AnythingOfType("web.GetJobsRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs?page=1", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r2.StatusCode)

		svc.On("GetJobs", mock.Anything, mock.AnythingOfType("web.GetJobsRequest")).
			Return(&web.JobPaginationResponse{Total: 1}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs?page=1", nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("GetRecruiterJobs validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recruiter/jobs?status=INVALID", nil))
		assert.Equal(t, fiber.StatusUnprocessableEntity, r1.StatusCode)

		svc.On("GetRecruiterJobs", mock.Anything, recruiterID, mock.AnythingOfType("web.GetRecruiterJobsRequest")).
			Return(nil, errs.ErrInternalServer).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recruiter/jobs?status=OPEN", nil))
		assert.Equal(t, fiber.StatusInternalServerError, r2.StatusCode)

		svc.On("GetRecruiterJobs", mock.Anything, recruiterID, mock.AnythingOfType("web.GetRecruiterJobsRequest")).
			Return(&web.JobPaginationResponse{Total: 1}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recruiter/jobs?status=OPEN", nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("GetJobById invalid UUID, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/not-a-uuid", nil))
		assert.Equal(t, fiber.StatusBadRequest, r1.StatusCode)

		svc.On("GetJobById", mock.Anything, jobID).Return(nil, errs.ErrJobNotFound).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusNotFound, r2.StatusCode)

		svc.On("GetJobById", mock.Anything, jobID).Return(&web.JobResponse{ID: jobID}, nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodGet, "/jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})

	t.Run("UpdateJob invalid UUID, validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPut, "/jobs/not-a-uuid", nil))
		assert.Equal(t, fiber.StatusBadRequest, r1.StatusCode)

		req2 := httptest.NewRequest(http.MethodPut, "/jobs/"+jobID.String(), bytes.NewBufferString(`{}`))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r2.StatusCode)

		payload, _ := json.Marshal(web.UpdateJobRequest{
			Title:          "Updated Title",
			Description:    "Updated Description",
			EmploymentType: "FULL_TIME",
			WorkMode:       "HYBRID",
			Location:       "Jakarta",
		})
		svc.On("UpdateJob", mock.Anything, recruiterID, jobID, mock.AnythingOfType("web.UpdateJobRequest")).
			Return(nil, errs.ErrJobForbidden).Once()
		req3 := httptest.NewRequest(http.MethodPut, "/jobs/"+jobID.String(), bytes.NewBuffer(payload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusForbidden, r3.StatusCode)

		svc.On("UpdateJob", mock.Anything, recruiterID, jobID, mock.AnythingOfType("web.UpdateJobRequest")).
			Return(&web.JobResponse{ID: jobID, Title: "Updated Title"}, nil).Once()
		req4 := httptest.NewRequest(http.MethodPut, "/jobs/"+jobID.String(), bytes.NewBuffer(payload))
		req4.Header.Set("Content-Type", "application/json")
		r4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusOK, r4.StatusCode)
	})

	t.Run("UpdateJobStatus invalid UUID, validation error, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodPatch, "/jobs/not-a-uuid/status", nil))
		assert.Equal(t, fiber.StatusBadRequest, r1.StatusCode)

		req2 := httptest.NewRequest(http.MethodPatch, "/jobs/"+jobID.String()+"/status", bytes.NewBufferString(`{"status":"INVALID"}`))
		req2.Header.Set("Content-Type", "application/json")
		r2, _ := app.Test(req2)
		assert.Equal(t, fiber.StatusUnprocessableEntity, r2.StatusCode)

		bodyPayload, _ := json.Marshal(map[string]any{
			"id":     jobID.String(),
			"status": "CLOSED",
		})
		svc.On("UpdateJobStatus", mock.Anything, recruiterID, jobID, "CLOSED").
			Return(nil, errs.ErrJobNotFound).Once()
		req3 := httptest.NewRequest(http.MethodPatch, "/jobs/"+jobID.String()+"/status", bytes.NewBuffer(bodyPayload))
		req3.Header.Set("Content-Type", "application/json")
		r3, _ := app.Test(req3)
		assert.Equal(t, fiber.StatusNotFound, r3.StatusCode)

		svc.On("UpdateJobStatus", mock.Anything, recruiterID, jobID, "CLOSED").
			Return(&web.JobResponse{ID: jobID, Status: domain.JobStatusClosed}, nil).Once()
		req4 := httptest.NewRequest(http.MethodPatch, "/jobs/"+jobID.String()+"/status", bytes.NewBuffer(bodyPayload))
		req4.Header.Set("Content-Type", "application/json")
		r4, _ := app.Test(req4)
		assert.Equal(t, fiber.StatusOK, r4.StatusCode)
	})

	t.Run("DeleteJob invalid UUID, service error, and success", func(t *testing.T) {
		r1, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/jobs/not-a-uuid", nil))
		assert.Equal(t, fiber.StatusBadRequest, r1.StatusCode)

		svc.On("DeleteJob", mock.Anything, recruiterID, jobID).Return(errs.ErrJobForbidden).Once()
		r2, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusForbidden, r2.StatusCode)

		svc.On("DeleteJob", mock.Anything, recruiterID, jobID).Return(nil).Once()
		r3, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/jobs/"+jobID.String(), nil))
		assert.Equal(t, fiber.StatusOK, r3.StatusCode)
	})
}
