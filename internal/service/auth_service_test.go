package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	req := web.RegisterRequest{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password123",
		Role:     domain.RoleCandidate,
	}

	t.Run("Hash error returns ErrInternalServer", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		tokenRepo := new(mockTokenRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)

		hasher.On("Hash", []byte(req.Password)).Return("", errors.New("bcrypt failure"))

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Register(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)
	})

	t.Run("User already exists returns ErrUserAlreadyExists", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		tokenRepo := new(mockTokenRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)

		hasher.On("Hash", []byte(req.Password)).Return("hashed-pw", nil)
		authRepo.On("IsUserExist", ctx, req.Email).Return(true)

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Register(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUserAlreadyExists, err)
	})

	t.Run("Repository register error returns error", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		tokenRepo := new(mockTokenRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)

		hasher.On("Hash", []byte(req.Password)).Return("hashed-pw", nil)
		authRepo.On("IsUserExist", ctx, req.Email).Return(false)
		authRepo.On("Register", ctx, mock.AnythingOfType("domain.User")).Return(nil, errors.New("db error"))

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Register(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Success returns RegisterResponse", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		tokenRepo := new(mockTokenRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)

		userID := uuid.New()
		hasher.On("Hash", []byte(req.Password)).Return("hashed-pw", nil)
		authRepo.On("IsUserExist", ctx, req.Email).Return(false)
		authRepo.On("Register", ctx, mock.AnythingOfType("domain.User")).Return(&domain.User{
			ID:    userID,
			Name:  req.Name,
			Email: req.Email,
			Role:  req.Role,
		}, nil)

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Register(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, userID, res.ID)
		assert.Equal(t, req.Name, res.Name)
		assert.Equal(t, req.Email, res.Email)
		assert.Equal(t, req.Role, res.Role)
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	req := web.LoginRequest{
		Email:    "jane@example.com",
		Password: "password123",
	}
	userID := uuid.New()
	foundUser := &domain.User{
		ID:    userID,
		Email: req.Email,
		Role:  domain.RoleCandidate,
		Auth: domain.Auth{
			Password: "hashed-password",
		},
	}

	t.Run("Email not found returns ErrInvalidCredentials", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		authRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))

		svc := NewAuthService(authRepo, new(mockPasswordHasher), new(mockJwtManager), new(mockTokenRepository))
		res, err := svc.Login(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidCredentials, err)
	})

	t.Run("Password mismatch returns ErrInvalidCredentials", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		hasher := new(mockPasswordHasher)

		authRepo.On("FindByEmail", ctx, req.Email).Return(foundUser, nil)
		hasher.On("Compare", []byte(foundUser.Auth.Password), []byte(req.Password)).Return(false)

		svc := NewAuthService(authRepo, hasher, new(mockJwtManager), new(mockTokenRepository))
		res, err := svc.Login(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidCredentials, err)
	})

	t.Run("GenerateTokenPair failure returns ErrInternalServer", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)

		authRepo.On("FindByEmail", ctx, req.Email).Return(foundUser, nil)
		hasher.On("Compare", []byte(foundUser.Auth.Password), []byte(req.Password)).Return(true)
		jwtMgr.On("GenerateTokenPair", userID, foundUser.Role).Return(nil, errors.New("jwt err"))

		svc := NewAuthService(authRepo, hasher, jwtMgr, new(mockTokenRepository))
		res, err := svc.Login(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)
	})

	t.Run("TokenRepository.Save failure returns ErrInternalServer", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		pair := &domain.TokenPair{
			AccessToken:      "acc",
			RefreshToken:     "ref",
			RefreshExpiredAt: time.Now().Add(time.Hour),
		}
		authRepo.On("FindByEmail", ctx, req.Email).Return(foundUser, nil)
		hasher.On("Compare", []byte(foundUser.Auth.Password), []byte(req.Password)).Return(true)
		jwtMgr.On("GenerateTokenPair", userID, foundUser.Role).Return(pair, nil)
		tokenRepo.On("Save", ctx, mock.AnythingOfType("domain.Token")).Return(nil, errors.New("save err"))

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Login(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)
	})

	t.Run("Success returns TokenPair", func(t *testing.T) {
		authRepo := new(mockAuthRepository)
		hasher := new(mockPasswordHasher)
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		pair := &domain.TokenPair{
			AccessToken:      "acc",
			RefreshToken:     "ref",
			RefreshExpiredAt: time.Now().Add(time.Hour),
		}
		authRepo.On("FindByEmail", ctx, req.Email).Return(foundUser, nil)
		hasher.On("Compare", []byte(foundUser.Auth.Password), []byte(req.Password)).Return(true)
		jwtMgr.On("GenerateTokenPair", userID, foundUser.Role).Return(pair, nil)
		tokenRepo.On("Save", ctx, mock.AnythingOfType("domain.Token")).Return(&domain.Token{}, nil)

		svc := NewAuthService(authRepo, hasher, jwtMgr, tokenRepo)
		res, err := svc.Login(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, pair, res)
	})
}

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()

	t.Run("RevokeToken error", func(t *testing.T) {
		tokenRepo := new(mockTokenRepository)
		tokenRepo.On("RevokeToken", ctx, "bad-token").Return(errors.New("revoke error"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), new(mockJwtManager), tokenRepo)
		err := svc.Logout(ctx, "bad-token")
		assert.EqualError(t, err, "revoke error")
	})

	t.Run("Success", func(t *testing.T) {
		tokenRepo := new(mockTokenRepository)
		tokenRepo.On("RevokeToken", ctx, "valid-token").Return(nil)

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), new(mockJwtManager), tokenRepo)
		err := svc.Logout(ctx, "valid-token")
		assert.NoError(t, err)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	req := web.RefreshTokenRequest{
		RefreshToken: "ref-token",
		AccessToken:  "acc-token",
	}

	t.Run("Invalid refresh token returns 401 AppError", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(nil, errors.New("Token has expired"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, new(mockTokenRepository))
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, "Token has expired", err.Error())
	})

	t.Run("Token not found in repository returns error", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(nil, errors.New("token not found"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "token not found")
	})

	t.Run("Revoked more than 10 seconds ago returns ErrTokenRevoked", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		revokedAt := time.Now().Add(-30 * time.Second)
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId:    userID,
			RevokedAt: &revokedAt,
			User:      domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrTokenRevoked, err)
	})

	t.Run("GenerateTokenPair error returns error", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId: userID,
			User:   domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)
		jwtMgr.On("GenerateTokenPair", userID, domain.RoleCandidate).Return(nil, errors.New("gen error"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "gen error")
	})

	t.Run("RevokeToken error returns error", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		newPair := &domain.TokenPair{AccessToken: "new-acc", RefreshToken: "new-ref"}
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId: userID,
			User:   domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)
		jwtMgr.On("GenerateTokenPair", userID, domain.RoleCandidate).Return(newPair, nil)
		tokenRepo.On("RevokeToken", ctx, req.RefreshToken).Return(errors.New("revoke fail"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "revoke fail")
	})

	t.Run("Save new token error returns error", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		newPair := &domain.TokenPair{AccessToken: "new-acc", RefreshToken: "new-ref"}
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId: userID,
			User:   domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)
		jwtMgr.On("GenerateTokenPair", userID, domain.RoleCandidate).Return(newPair, nil)
		tokenRepo.On("RevokeToken", ctx, req.RefreshToken).Return(nil)
		tokenRepo.On("Save", ctx, mock.AnythingOfType("domain.Token")).Return(nil, errors.New("save fail"))

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "save fail")
	})

	t.Run("Success when token not yet revoked", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		newPair := &domain.TokenPair{AccessToken: "new-acc", RefreshToken: "new-ref"}
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId: userID,
			User:   domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)
		jwtMgr.On("GenerateTokenPair", userID, domain.RoleCandidate).Return(newPair, nil)
		tokenRepo.On("RevokeToken", ctx, req.RefreshToken).Return(nil)
		tokenRepo.On("Save", ctx, mock.AnythingOfType("domain.Token")).Return(&domain.Token{}, nil)

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, newPair, res.TokenPair)
	})

	t.Run("Success within 10s grace period of recently revoked token", func(t *testing.T) {
		jwtMgr := new(mockJwtManager)
		tokenRepo := new(mockTokenRepository)

		recentRevoke := time.Now().Add(-2 * time.Second)
		newPair := &domain.TokenPair{AccessToken: "new-acc", RefreshToken: "new-ref"}
		jwtMgr.On("ValidateRefreshToken", req.RefreshToken).Return(&domain.JwtCustomClaims{UserID: userID}, nil)
		tokenRepo.On("FindTokenWithUser", ctx, userID, req.RefreshToken).Return(&domain.Token{
			UserId:    userID,
			RevokedAt: &recentRevoke,
			User:      domain.User{ID: userID, Role: domain.RoleCandidate},
		}, nil)
		jwtMgr.On("GenerateTokenPair", userID, domain.RoleCandidate).Return(newPair, nil)
		tokenRepo.On("Save", ctx, mock.AnythingOfType("domain.Token")).Return(&domain.Token{}, nil)

		svc := NewAuthService(new(mockAuthRepository), new(mockPasswordHasher), jwtMgr, tokenRepo)
		res, err := svc.RefreshToken(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, newPair, res.TokenPair)
	})
}
