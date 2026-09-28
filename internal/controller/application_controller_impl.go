package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/request"
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/errs"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type applicationControllerImpl struct {
	applicationService service.ApplicationService
}

func NewApplicationController(applicationService service.ApplicationService) ApplicationController {
	return &applicationControllerImpl{
		applicationService: applicationService,
	}
}

// POST /api/jobs/:id/applications (Candidate Only)
func (ctrl *applicationControllerImpl) ApplyJob(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.ApplyJobRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}

	res, err := ctrl.applicationService.ApplyJob(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Apply job success", res)
}

// GET /api/jobs/:id/applications (Recruiter Only)
func (ctrl *applicationControllerImpl) GetJobApplications(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.GetJobApplicationsRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}

	res, err := ctrl.applicationService.GetJobApplications(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Get job applications success", res)
}

// GET /api/candidate/applications (Candidate Only)
func (ctrl *applicationControllerImpl) GetCandidateApplications(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.GetCandidateApplicationsRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	res, err := ctrl.applicationService.GetCandidateApplications(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Get candidate applications success", res)
}

// GET /api/applications/:id (Candidate or Recruiter)
func (ctrl *applicationControllerImpl) GetApplicationById(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	appIdStr := c.Params("id")
	appId, err := uuid.Parse(appIdStr)
	if err != nil {
		return errs.ErrInvalidInput
	}

	res, err := ctrl.applicationService.GetApplicationById(ctx, session.UserId, string(session.Role), appId)
	if err != nil {
		return err
	}

	return response.OK(c, "Get application detail success", res)
}

// PATCH /api/applications/:id/status (Recruiter Only)
func (ctrl *applicationControllerImpl) UpdateStatus(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.UpdateApplicationStatusRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}

	res, err := ctrl.applicationService.UpdateStatus(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Update application status success", res)
}

// PATCH /api/candidate/applications/:id/withdraw (Candidate Only)
func (ctrl *applicationControllerImpl) WithdrawApplication(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.WithdrawApplicationRequest
	err = c.Bind().All(&req)
	if err != nil {
		return err
	}

	res, err := ctrl.applicationService.WithdrawApplication(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Withdraw application success", res)
}
