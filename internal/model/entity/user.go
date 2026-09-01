package model

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "RECRUITER"
	RoleUser  UserRole = "CANDIDATE"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email;unique"`
	Password  string    `gorm:"column:password"`
	Role      UserRole  `gorm:"column:role;default:'CANDIDATE'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateDate;autoUpdateDate"`
}
