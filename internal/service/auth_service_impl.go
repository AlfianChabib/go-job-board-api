package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
	}
}

func (service *AuthServiceImpl) Register(ctx fiber.Ctx, data web.RegisterRequest) (*domain.User, error) {
	newHash, err := utils.HashPassword([]byte(data.Password))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	existUser := service.AuthRepository.IsUserAlreadyExist(ctx, data.Email)
	if existUser {
		return nil, fiber.NewError(fiber.StatusConflict, "User already exist")
	}

	user, err := service.AuthRepository.Register(ctx, domain.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(newHash),
	})
	return user, err
}

func (a *AuthServiceImpl) Login(ctx fiber.Ctx, data web.LoginRrequest) error {
	panic("TODO: Implement")
}

func (a *AuthServiceImpl) Logout(ctx fiber.Ctx) error {
	panic("TODO: Implement")
}
