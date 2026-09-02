package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"

	"github.com/gofiber/fiber/v3"
)

type authServiceImpl struct {
	AuthRepository repository.AuthRepository
	PasswordHasher domain.PasswordHasher
}

func NewAuthService(authRepository repository.AuthRepository, passwordHasher domain.PasswordHasher) AuthService {
	return &authServiceImpl{
		AuthRepository: authRepository,
		PasswordHasher: passwordHasher,
	}
}

func (service *authServiceImpl) Register(ctx fiber.Ctx, data web.RegisterRequest) (*web.RegisterResponse, error) {
	newHash, err := service.PasswordHasher.Hash([]byte(data.Password))
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
		Password: newHash,
	})

	return &web.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, err
}

func (a *authServiceImpl) Login(ctx fiber.Ctx, data web.LoginRrequest) error {
	panic("TODO: Implement")
}

func (a *authServiceImpl) Logout(ctx fiber.Ctx) error {
	panic("TODO: Implement")
}
