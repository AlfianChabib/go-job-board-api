package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type candidateController struct {
	candidateService service.CandidateService
}

func NewCandidateController(candidateService service.CandidateService) CandidateController {
	return &candidateController{
		candidateService: candidateService,
	}
}

func (controller *candidateController) Get(c fiber.Ctx) error {
	ctx := c.Context()
	userId := c.Locals("userId").(uuid.UUID)

	if userId == uuid.Nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	candidate, err := controller.candidateService.Get(ctx, userId)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get candidate", candidate)
}

func (controller *candidateController) Update(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.UpdateCandidateRequest
	userId := c.Locals("userId").(uuid.UUID)

	if userId == uuid.Nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	candidate, err := controller.candidateService.Update(ctx, domain.Profile{
		UserId:   userId,
		Headline: &req.Headline,
		Phone:    &req.Phone,
	})
	if err != nil {
		return err
	}

	return response.OK(c, "Update candidate profile success", candidate)
}
