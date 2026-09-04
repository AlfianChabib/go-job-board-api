package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"

	"github.com/gofiber/fiber/v3"
)

type authServiceImpl struct {
	AuthRepository  repository.AuthRepository
	TokenRepository repository.TokenRepository
	PasswordHasher  domain.PasswordHasher
	JwtManager      domain.JwtManager
}

func NewAuthService(
	authRepository repository.AuthRepository,
	passwordHasher domain.PasswordHasher,
	jwtManager domain.JwtManager,
	tokenRepository repository.TokenRepository,
) AuthService {
	return &authServiceImpl{
		AuthRepository:  authRepository,
		PasswordHasher:  passwordHasher,
		JwtManager:      jwtManager,
		TokenRepository: tokenRepository,
	}
}

func (service *authServiceImpl) Register(ctx context.Context, data web.RegisterRequest) (*web.RegisterResponse, error) {
	newHash, err := service.PasswordHasher.Hash([]byte(data.Password))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	if service.AuthRepository.IsUserExist(ctx, data.Email) {
		return nil, fiber.NewError(fiber.StatusConflict, "User already exist")
	}

	user, err := service.AuthRepository.Register(ctx, domain.User{
		Name:  data.Name,
		Email: data.Email,
		Auth: domain.Auth{
			Email:    data.Email,
			Password: newHash,
		},
	})
	if err != nil {
		return nil, err
	}

	return &web.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (service *authServiceImpl) Login(ctx context.Context, data web.LoginRrequest) (*domain.TokenPair, error) {
	user, err := service.AuthRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Wrong Password or Email")
	}

	matchedPassword := service.PasswordHasher.Compare([]byte(user.Auth.Password), []byte(data.Password))
	if !matchedPassword {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Wrong Password or Email")
	}

	tokens, err := service.JwtManager.GenerateTokenPair(user.ID.String(), user.Role)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	_, err = service.TokenRepository.Save(ctx, domain.Token{
		UserId:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.RefreshExpiredAt,
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return tokens, nil
}

func (a *authServiceImpl) Logout(ctx context.Context) error {
	panic("TODO: Implement")
}
