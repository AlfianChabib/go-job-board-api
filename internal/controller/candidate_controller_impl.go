package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper"
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"path/filepath"
	"slices"
	"strings"

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

func (controller *candidateController) UpdateAvatar(c fiber.Ctx) error {
	maxFileSize := 5 * 1024 * 1024
	allowedFileContentTypes := []string{"image/jpg", "image/jpeg", "image/png"}
	allowedFileFormats := []string{".jpg", "jpeg", ".png"}

	ctx := c.Context()
	userId := c.Locals("userId").(uuid.UUID)

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file not found")
	}

	if fileHeader.Size > int64(maxFileSize) {
		return fiber.NewError(fiber.StatusBadRequest, "file to large")
	}

	if err = helper.ValidateFilename(fileHeader.Filename); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedFileFormat := slices.Contains(allowedFileFormats, fileExtension)
	if !allowedFileFormat {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	fileContentType := fileHeader.Header.Get("Content-Type")
	allowedFileContentType := slices.Contains(allowedFileContentTypes, fileContentType)
	if !allowedFileContentType {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process file")
	}
	defer file.Close()

	req := web.UpdateCandidateAvatarRequest{
		UserId:     userId,
		File:       file,
		FileSize:   fileHeader.Size,
		ContenType: fileContentType,
		Extension:  fileExtension,
	}

	responseUrl, err := controller.candidateService.UploadAvatar(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server err")
	}

	return response.OK(c, "Success upload avatar", responseUrl)
}
