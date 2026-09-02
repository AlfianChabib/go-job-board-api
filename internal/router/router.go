package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/database"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

func Initialize(router *fiber.App) {
	db := database.OpenConnection()
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
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
