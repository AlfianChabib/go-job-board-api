package router

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/database"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func InitializeRouter(router *fiber.App, env *config.Env) {
	db := database.OpenConnection()
	authRepository := repository.NewAuthRepository(db)
	tokenRepository := repository.NewTokenRepository(db)

	passwordHasher := utils.NewBcryptHasher(bcrypt.DefaultCost)
	jwtManager := utils.NewJwtManager(env.AccessSecretKey, env.AccessDuration, env.RefreshSecretKey, env.RefreshDuration)

	authService := service.NewAuthService(authRepository, passwordHasher, jwtManager, tokenRepository)
	authController := controller.NewAuthController(authService, env.AppEnv)

	router.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.Post("/register", authController.Register)
		auth.Post("/login", authController.Login)
		auth.Post("/logout", authController.LogOut)
	}
}
