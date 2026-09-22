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

	"github.com/minio/minio-go/v7"

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

func ProvideRecruiterController(recruiterService service.RecruiterService) controller.RecruiterController {
	return controller.NewRecruiterController(recruiterService)
}

func ProvideCompanyController(companyService service.CompanyService) controller.CompanyController {
	return controller.NewCompanyController(companyService)
}

func ProvideStoragerepository(store *minio.Client, env *config.Env) repository.StorageRepository {
	return repository.NewStorageRepository(store, env.MinioPublicUrl, env.MinioAvatarBucket, env.MinioCvBucket)
}

var authSet = wire.NewSet(
	repository.NewAuthRepository,
	repository.NewTokenRepository,
	service.NewAuthService,
	ProvideAuthController,
)

var candidateSet = wire.NewSet(
	repository.NewCandidateRepository,
	ProvideStoragerepository,
	service.NewCandidateService,
	ProvideCandidateController,
)

var recruiterSet = wire.NewSet(
	service.NewRecruiterService,
	ProvideRecruiterController,
)

var companySet = wire.NewSet(
	repository.NewCompanyRepository,
	service.NewCompanyService,
	ProvideCompanyController,
)

func InitializeApp(env *config.Env) (*fiber.App, error) {
	wire.Build(
		database.OpenConnection,
		database.OpenMinioClient,
		validator.NewValidator,
		ProvideBcryptHasher,
		ProvideJwtManager,
		middleware.NewMiddleware,
		authSet,
		candidateSet,
		recruiterSet,
		companySet,
		NewApp,
	)
	return nil, nil
}
