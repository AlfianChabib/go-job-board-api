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
) {
	router.Get("/", func(c fiber.Ctx) error {
		return response.Message(c, fiber.StatusOK, "Hello world!")
	})

	api := router.Group("/api")
	SetupAuthRoutes(api, authController)
	SetupCandidateRoutes(api, middleware, candidateController)
	SetupRecruiterRoutes(api, middleware, recruiterController)
	SetupUserRoutes(api)
}
