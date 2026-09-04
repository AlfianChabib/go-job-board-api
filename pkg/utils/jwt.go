package utils

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct {
	AccessSecretKey  string
	AccessDuration   int
	RefreshSecretKey string
	RefreshDuration  int
}

func NewJwtManager(
	accessSecretKey string,
	accessDuration int,
	refreshSecretKey string,
	refreshDuration int,
) domain.JwtManager {
	return &jwtManager{
		AccessSecretKey:  accessSecretKey,
		AccessDuration:   accessDuration,
		RefreshSecretKey: refreshSecretKey,
		RefreshDuration:  refreshDuration,
	}
}

func (j *jwtManager) GenerateTokenPair(userID string, role domain.UserRole) (*domain.TokenPair, error) {
	accessExpiredAt := time.Now().Add(time.Duration(j.AccessDuration * int(time.Second)))
	refreshExpiredAt := time.Now().Add(time.Duration(j.RefreshDuration * int(time.Second)))

	accessClaims := &domain.JwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiredAt),
		},
	}
	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.AccessSecretKey))

	refreshClaims := &domain.JwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiredAt),
		},
	}
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.RefreshSecretKey))

	return &domain.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiredAt:  accessExpiredAt,
		RefreshExpiredAt: refreshExpiredAt,
	}, nil
}

func (j *jwtManager) ValidateToken(tokenString string) (*domain.JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JwtCustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing token")
		}
		return []byte(j.AccessSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*domain.JwtCustomClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
