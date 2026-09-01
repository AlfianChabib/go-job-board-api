package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

func Initialize(router *fiber.App) {
	authService := service.NewAuthService()
	authController := controller.NewAuthController(authService)

	router.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.Post("/register", authController.Register)
	}
}
