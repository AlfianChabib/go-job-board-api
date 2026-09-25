package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
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

func (d *dataController) GetSkills(c fiber.Ctx) error {
	ctx := c.Context()
	skills, err := d.dataService.GetSkills(ctx)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get skills", skills)
}

func (d *dataController) GetCurrencyCodes(c fiber.Ctx) error {
	ctx := c.Context()
	currencies, err := d.dataService.GetCurrencyCodes(ctx)
	if err != nil {
		return err
	}

	return response.OK(c, "Success get currency codes", currencies)
}
