package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/errs"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type companyControllerImpl struct {
	companyService service.CompanyService
}

func NewCompanyController(companyService service.CompanyService) CompanyController {
	return &companyControllerImpl{
		companyService: companyService,
	}
}

func (ctrl *companyControllerImpl) GetCompanies(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.GetCompaniesRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	res, err := ctrl.companyService.GetCompanies(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get companies", res)
}

func (ctrl *companyControllerImpl) GetCompanyById(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.GetCompanyByIdRequest
	if err := c.Bind().URI(&req); err != nil {
		return err
	}

	if req.ID == uuid.Nil {
		companyIdStr := c.Params("id")
		if companyIdStr == "" {
			companyIdStr = c.Params("companyId")
		}
		parsed, err := uuid.Parse(companyIdStr)
		if err != nil {
			return errs.ErrInvalidInput
		}
		req.ID = parsed
	}

	res, err := ctrl.companyService.GetCompanyById(ctx, req.ID)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get company detail", res)
}
