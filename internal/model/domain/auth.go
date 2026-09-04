package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Auth struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	UserId    uuid.UUID `gorm:"column:user_id;type:uuid"`
	Password  string    `gorm:"column:password;type:varchar;size:255"`
	Email     string    `gorm:"column:email;type:varchar;size:255"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateDate;autoUpdateDate"`
	User      *User     `gorm:"foreignKey:UserId;references:ID"`
}

func (a *Auth) TableName() string {
	return "auths"
}

func (a *Auth) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		a.ID = uuidV7
	}

	return
}

type PasswordHasher interface {
	Hash(password []byte) (string, error)
	Compare(hashedPassword, password []byte) bool
}
