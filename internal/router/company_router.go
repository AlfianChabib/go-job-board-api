package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupCompanyRoutes(
	router fiber.Router,
	companyController controller.CompanyController,
) {
	companies := router.Group("/companies")

	companies.Get("/", companyController.GetCompanies)
	companies.Get("/:id", companyController.GetCompanyById)
}
