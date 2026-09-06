package router

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func InitializeRoutes(
	router *fiber.App,
	env *config.Env,
	authController controller.AuthController,
) {
	router.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Hello World",
		})
	})

	api := router.Group("/api")
	SetupAuthRoutes(api, authController)
	SetupUserRoutes(api)
}
