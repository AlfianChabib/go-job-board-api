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

	// repository
	authRepository := repository.NewAuthRepository(db)
	tokenRepository := repository.NewTokenRepository(db)

	// utils
	passwordHasher := utils.NewBcryptHasher(bcrypt.DefaultCost)
	jwtManager := utils.NewJwtManager(env.AccessSecretKey, env.AccessDuration, env.RefreshSecretKey, env.RefreshDuration)

	// service
	authService := service.NewAuthService(authRepository, passwordHasher, jwtManager, tokenRepository)
	authController := controller.NewAuthController(authService, env.AppEnv)

	// middleware
	// protected := middleware.Protected(jwtManager)

	router.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := router.Group("/api")
	SetupAuthRoutes(api, authController)
}
