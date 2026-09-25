package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupDataRoutes(
	router fiber.Router,
	dataController controller.DataController,
) {
	router.Get("/currencies", dataController.GetCurrencyCodes)
	router.Get("/skills", dataController.GetSkills)
}
