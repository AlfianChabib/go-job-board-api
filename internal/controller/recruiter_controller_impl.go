package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper"
	"AlfianChabib/go-job-board-api/internal/helper/request"
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/utils"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type recruiterController struct {
	recruiterService service.RecruiterService
}

func NewRecruiterController(recruiterService service.RecruiterService) RecruiterController {
	return &recruiterController{
		recruiterService: recruiterService,
	}
}

func (ctrl *recruiterController) CreateCompany(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.CreateCompanyRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	req.Description = utils.NilIfEmpty(req.Description)
	req.Website = utils.NilIfEmpty(req.Website)

	res, err := ctrl.recruiterService.CreateCompany(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Create company success", res)
}

func (ctrl *recruiterController) GetCompany(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	res, err := ctrl.recruiterService.GetCompany(ctx, session.UserId)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get company", res)
}

func (ctrl *recruiterController) UpdateCompany(c fiber.Ctx) error {
	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	var req web.UpdateCompanyRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	req.Website = utils.NilIfEmpty(req.Website)
	req.Description = utils.NilIfEmpty(req.Description)

	res, err := ctrl.recruiterService.UpdateCompany(ctx, session.UserId, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success update company", res)
}

func (ctrl *recruiterController) UpdateLogo(c fiber.Ctx) error {
	maxFileSize := 2 * 1024 * 1024 // 2MB
	allowedFileContentTypes := []string{"image/png", "image/jpeg", "image/webp"}
	allowedFileFormats := []string{".png", ".jpg", ".jpeg", ".webp"}

	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fileHeader, err = c.FormFile("logo")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file not found")
		}
	}

	if fileHeader.Size > int64(maxFileSize) {
		return fiber.NewError(fiber.StatusBadRequest, "file too large")
	}

	if err = helper.ValidateFilename(fileHeader.Filename); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !slices.Contains(allowedFileFormats, fileExtension) {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	fileContentType := fileHeader.Header.Get("Content-Type")
	if !slices.Contains(allowedFileContentTypes, fileContentType) {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process file")
	}
	defer file.Close()

	req := web.UpdateCompanyLogoRequest{
		RecruiterId: session.UserId,
		File:        file,
		FileSize:    fileHeader.Size,
		ContentType: fileContentType,
		Extension:   fileExtension,
	}

	res, err := ctrl.recruiterService.UploadLogo(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success upload company logo", res)
}

func (ctrl *recruiterController) UpdateBanner(c fiber.Ctx) error {
	maxFileSize := 5 * 1024 * 1024 // 5MB
	allowedFileContentTypes := []string{"image/png", "image/jpeg", "image/webp"}
	allowedFileFormats := []string{".png", ".jpg", ".jpeg", ".webp"}

	ctx := c.Context()
	session, err := request.GetLocalSession(c)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fileHeader, err = c.FormFile("banner")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file not found")
		}
	}

	if fileHeader.Size > int64(maxFileSize) {
		return fiber.NewError(fiber.StatusBadRequest, "file too large")
	}

	if err = helper.ValidateFilename(fileHeader.Filename); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !slices.Contains(allowedFileFormats, fileExtension) {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	fileContentType := fileHeader.Header.Get("Content-Type")
	if !slices.Contains(allowedFileContentTypes, fileContentType) {
		return fiber.NewError(fiber.StatusBadRequest, "file format not allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process file")
	}
	defer file.Close()

	req := web.UpdateCompanyBannerRequest{
		RecruiterId: session.UserId,
		File:        file,
		FileSize:    fileHeader.Size,
		ContentType: fileContentType,
		Extension:   fileExtension,
	}

	res, err := ctrl.recruiterService.UploadBanner(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success upload company banner", res)
}
