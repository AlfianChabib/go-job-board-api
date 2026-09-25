package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type dataController struct {
	dataService service.DataService
}

func NewDataController(dataService service.DataService) DataController {
	return &dataController{
		dataService: dataService,
	}
}

func (ctrl *dataController) GetSkills(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.GetDataRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	skills, err := ctrl.dataService.GetSkills(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get skills", skills)
}

func (ctrl *dataController) GetCurrencyCodes(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.GetDataRequest
	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	currencies, err := ctrl.dataService.GetCurrencyCodes(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get currency codes", currencies)
}
