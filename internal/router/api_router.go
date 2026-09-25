package router

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func InitializeRoutes(
	router *fiber.App,
	env *config.Env,
	middleware middleware.Middleware,
	authController controller.AuthController,
	candidateController controller.CandidateController,
	recruiterController controller.RecruiterController,
	companyController controller.CompanyController,
	dataController controller.DataController,
) {
	router.Get("/", func(c fiber.Ctx) error {
		return response.Message(c, fiber.StatusOK, "Hello world!")
	})

	api := router.Group("/api")
	SetupDataRoutes(api, dataController)
	SetupAuthRoutes(api, authController)
	SetupCandidateRoutes(api, middleware, candidateController)
	SetupRecruiterRoutes(api, middleware, recruiterController)
	SetupCompanyRoutes(api, companyController)
	SetupUserRoutes(api)
}
