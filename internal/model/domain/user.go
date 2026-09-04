package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "RECRUITER"
	RoleUser  UserRole = "CANDIDATE"
)

type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email;unique"`
	Role      UserRole  `gorm:"column:role;default:'CANDIDATE'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateDate;autoUpdateDate"`
	Auth      Auth      `gorm:"foreignKey:UserId;references:ID"`
	Tokens    []Token   `gorm:"foreignKey:UserId;references:ID"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		u.ID = uuidV7
	}

	return
}

type JwtCustomClaims struct {
	UserID string   `json:"user_id"`
	Role   UserRole `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

type JwtManager interface {
	GenerateTokenPair(userID string, role UserRole) (*TokenPair, error)
	ValidateToken(tokenString string) (*JwtCustomClaims, error)
}
