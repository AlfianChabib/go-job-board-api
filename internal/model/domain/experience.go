package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Experience struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	ProfileId   uuid.UUID  `gorm:"column:profile_id;type:uuid;not null"`
	CompanyName string     `gorm:"column:company_name;type:varchar;size:255;not null"`
	Position    string     `gorm:"column:position;type:varchar;size:255;not null"`
	StartDate   time.Time  `gorm:"column:start_date;type:date;not null"`
	EndDate     *time.Time `gorm:"column:end_date;type:date"`
	IsCurrent   bool       `gorm:"column:is_current;type:bool;default:false"`
	Description *string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	Profile     *Profile   `gorm:"foreignKey:ProfileId;references:ID"`
}

func (e *Experience) TableName() string {
	return "experiences"
}

func (e *Experience) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		e.ID = uuidV7
	}
	return nil
}
