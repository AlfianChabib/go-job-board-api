package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid"`
	UserId       uuid.UUID `gorm:"column:user_id;type:uuid"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar;size:255"`
	IsRevoked    bool      `gorm:"column:is_revoked;type:bool;default:false"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt    time.Time `gorm:"column:created_at;autoCreateDate;autoUpdateDate"`
	User         User      `gorm:"foreignKey:UserId;references:ID"`
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
