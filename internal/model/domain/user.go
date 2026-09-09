package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleRecruiter UserRole = "RECRUITER"
	RoleCandidate UserRole = "CANDIDATE"
)

type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email;unique"`
	Role      UserRole  `gorm:"column:role;default:'CANDIDATE'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	Auth      Auth      `gorm:"foreignKey:UserId;references:ID"`
	Tokens    []Token   `gorm:"foreignKey:UserId;references:ID"`
	Profile   *Profile  `gorm:"foreignKey:UserId;references:ID"`
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
