package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeSize string

const (
	EmployeeSize1To50       EmployeeSize = "1-50"
	EmployeeSize51To200     EmployeeSize = "51-200"
	EmployeeSize201To500    EmployeeSize = "201-500"
	EmployeeSizeMoreThan500 EmployeeSize = "500+"
)

type Company struct {
	ID           uuid.UUID    `gorm:"column:id;type:uuid;primaryKey"`
	RecruiterId  uuid.UUID    `gorm:"column:recruiter_id;type:uuid;unique;not null"`
	Name         string       `gorm:"column:name;type:varchar;size:255;not null"`
	LogoUrl      *string      `gorm:"column:logo_url;type:varchar;size:255"`
	BannerUrl    *string      `gorm:"column:banner_url;type:varchar;size:255"`
	Website      *string      `gorm:"column:website;type:varchar;size:255"`
	Industry     string       `gorm:"column:industry;type:varchar;size:255;not null"`
	EmployeeSize EmployeeSize `gorm:"column:employee_size;type:varchar;size:50;not null"`
	Description  *string      `gorm:"column:description;type:text"`
	Location     string       `gorm:"column:location;type:varchar;size:255;not null"`
	CreatedAt    time.Time    `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt    time.Time    `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	Recruiter    *User        `gorm:"foreignKey:RecruiterId;references:ID"`
}

func (c *Company) TableName() string {
	return "companies"
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		c.ID = uuidV7
	}
	return nil
}
