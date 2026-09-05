package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupAuthRoutes(router fiber.Router, authController controller.AuthController) {
	auth := router.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.LogOut)
	auth.Post("/refresh", authController.RefreshToken)
}
