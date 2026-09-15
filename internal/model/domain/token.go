package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID           uuid.UUID  `gorm:"column:id;type:uuid"`
	UserId       uuid.UUID  `gorm:"column:user_id;type:uuid"`
	RefreshToken string     `gorm:"column:refresh_token;type:varchar;size:255"`
	RevokedAt    *time.Time `gorm:"column:revoked_at"`
	ExpiresAt    time.Time  `gorm:"column:expires_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt    time.Time  `gorm:"column:created_at;autoCreateDate;autoUpdateDate"`
	User         User       `gorm:"foreignKey:UserId;references:ID"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		t.ID = uuidV7
	}

	return nil
}

type JwtCustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   UserRole  `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

type JwtManager interface {
	GenerateTokenPair(userID uuid.UUID, role UserRole) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (*JwtCustomClaims, error)
	ValidateRefreshToken(tokenString string) (*JwtCustomClaims, error)
}
