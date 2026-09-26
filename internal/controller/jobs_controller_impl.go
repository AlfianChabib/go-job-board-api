package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/request"
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"AlfianChabib/go-job-board-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type jobControllerImpl struct {
	jobService service.JobService
}

func NewJobController(jobService service.JobService) JobController {
	return &jobControllerImpl{
		jobService: jobService,
	}
}

// POST /api/jobs (Recruiter Only)
func (ctrl *jobControllerImpl) CreateJob(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.CreateJobRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	res, err := ctrl.jobService.CreateJob(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Create job success", res)
}

// GET /api/jobs (Public)
func (ctrl *jobControllerImpl) GetJobs(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.GetJobsRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	req.MinSalary = utils.NilIfZero(req.MinSalary)

	res, err := ctrl.jobService.GetJobs(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Get jobs success", res)
}

// GET /api/recruiter/jobs (Recruiter Only)
func (ctrl *jobControllerImpl) GetRecruiterJobs(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.GetRecruiterJobsRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	res, err := ctrl.jobService.GetRecruiterJobs(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Get recruiter jobs success", res)
}

// GET /api/jobs/:id (Public)
func (ctrl *jobControllerImpl) GetJobById(c fiber.Ctx) error {
	ctx := c.Context()
	jobIdStr := c.Params("id")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		return errs.ErrInvalidInput
	}

	res, err := ctrl.jobService.GetJobById(ctx, jobId)
	if err != nil {
		return err
	}

	return response.OK(c, "Get job detail success", res)
}

// PUT /api/jobs/:id (Recruiter Only)
func (ctrl *jobControllerImpl) UpdateJob(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	jobIdStr := c.Params("id")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		return errs.ErrInvalidInput
	}

	var req web.UpdateJobRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}
	req.ID = jobId

	res, err := ctrl.jobService.UpdateJob(ctx, session.UserId, jobId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Update job success", res)
}

// PATCH /api/jobs/:id/status (Recruiter Only)
func (ctrl *jobControllerImpl) UpdateJobStatus(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	jobIdStr := c.Params("id")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		return errs.ErrInvalidInput
	}

	var req web.UpdateJobStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	req.ID = jobId

	res, err := ctrl.jobService.UpdateJobStatus(ctx, session.UserId, jobId, req.Status)
	if err != nil {
		return err
	}

	return response.OK(c, "Update job status success", res)
}

// DELETE /api/jobs/:id (Recruiter Only)
func (ctrl *jobControllerImpl) DeleteJob(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	jobIdStr := c.Params("id")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		return errs.ErrInvalidInput
	}

	if err := ctrl.jobService.DeleteJob(ctx, session.UserId, jobId); err != nil {
		return err
	}

	return response.Message(c, fiber.StatusOK, "Delete job success")
}
