package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"net/http"
	"time"
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
		return nil, errs.ErrInternalServer
	}

	if service.AuthRepository.IsUserExist(ctx, data.Email) {
		return nil, errs.ErrUserAlreadyExists
	}

	user, err := service.AuthRepository.Register(ctx, domain.User{
		Name:  data.Name,
		Email: data.Email,
		Role:  data.Role,
		Auth: domain.Auth{
			Email:    data.Email,
			Password: newHash,
		},
		Profile: &domain.Profile{},
	})
	if err != nil {
		return nil, err
	}

	return &web.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (service *authServiceImpl) Login(ctx context.Context, data web.LoginRequest) (*domain.TokenPair, error) {
	user, err := service.AuthRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	matchedPassword := service.PasswordHasher.Compare([]byte(user.Auth.Password), []byte(data.Password))
	if !matchedPassword {
		return nil, errs.ErrInvalidCredentials
	}

	tokens, err := service.JwtManager.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	_, err = service.TokenRepository.Save(ctx, domain.Token{
		UserId:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.RefreshExpiredAt,
	})
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	return tokens, nil
}

func (service *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	err := service.TokenRepository.RevokeToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (service *authServiceImpl) RefreshToken(ctx context.Context, data web.RefreshTokenRequest) (*web.RefreshTokenResponse, error) {
	decodedToken, err := service.JwtManager.ValidateRefreshToken(data.RefreshToken)
	if err != nil {
		return nil, errs.New(http.StatusUnauthorized, err.Error())
	}

	userWithToken, err := service.TokenRepository.FindTokenWithUser(ctx, decodedToken.UserID, data.RefreshToken)
	if err != nil {
		return nil, err
	}

	if userWithToken.RevokedAt != nil {
		timeSinceRevoked := time.Since(*userWithToken.RevokedAt)
		if timeSinceRevoked > (10 * time.Second) {
			return nil, errs.ErrTokenRevoked
		}
	}

	newTokens, err := service.JwtManager.GenerateTokenPair(userWithToken.User.ID, userWithToken.User.Role)
	if err != nil {
		return nil, err
	}

	if userWithToken.RevokedAt == nil {
		err = service.TokenRepository.RevokeToken(ctx, data.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	_, err = service.TokenRepository.Save(ctx, domain.Token{
		UserId:       userWithToken.User.ID,
		RefreshToken: newTokens.RefreshToken,
		ExpiresAt:    newTokens.RefreshExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	return &web.RefreshTokenResponse{TokenPair: newTokens}, nil
}
