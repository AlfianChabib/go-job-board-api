//go:build wireinject
// +build wireinject

package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/database"
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/internal/service"
	"AlfianChabib/go-job-board-api/pkg/utils"
	"AlfianChabib/go-job-board-api/pkg/validator"

	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
)

func ProvideBcryptHasher() domain.PasswordHasher {
	return utils.NewBcryptHasher(bcrypt.DefaultCost)
}

func ProvideJwtManager(env *config.Env) domain.JwtManager {
	return utils.NewJwtManager(
		env.AccessSecretKey,
		env.AccessDuration,
		env.RefreshSecretKey,
		env.RefreshDuration,
	)
}

func ProvideAuthController(authService service.AuthService, env *config.Env) controller.AuthController {
	return controller.NewAuthController(authService, env.AppEnv)
}

func ProvideCandidateController(candidateService service.CandidateService) controller.CandidateController {
	return controller.NewCandidateController(candidateService)
}

var authSet = wire.NewSet(
	repository.NewAuthRepository,
	repository.NewTokenRepository,
	service.NewAuthService,
	ProvideAuthController,
)

var candidateSet = wire.NewSet(
	repository.NewCandidateRepository,
	service.NewCandidateService,
	ProvideCandidateController,
)

func InitializeApp(env *config.Env) (*fiber.App, error) {
	wire.Build(
		database.OpenConnection,
		validator.NewValidator,
		ProvideBcryptHasher,
		ProvideJwtManager,
		middleware.NewMiddleware,
		authSet,
		candidateSet,
		NewApp,
	)
	return nil, nil
}
