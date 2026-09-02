package domain

import "github.com/golang-jwt/jwt/v5"

type PasswordHasher interface {
	Hash(password []byte) (string, error)
	Compare(hashedPassword, password []byte) bool
}

type JwtCustomClaims struct {
	UserID string   `json:"user_id"`
	Role   UserRole `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtManager interface {
	GenerateTokenPair(userID string, role UserRole) (*TokenPair, error)
	ValidateToken(tokenString string) (*JwtCustomClaims, error)
}
