package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmploymentType string

const (
	EmploymentTypeFullTime   EmploymentType = "FULL_TIME"
	EmploymentTypePartTime   EmploymentType = "PART_TIME"
	EmploymentTypeContract   EmploymentType = "CONTRACT"
	EmploymentTypeInternship EmploymentType = "INTERNSHIP"
	EmploymentTypeFreelance  EmploymentType = "FREELANCE"
)

type WorkMode string

const (
	WorkModeOnSite WorkMode = "ON_SITE"
	WorkModeHybrid WorkMode = "HYBRID"
	WorkModeRemote WorkMode = "REMOTE"
)

type JobStatus string

const (
	JobStatusOpen   JobStatus = "OPEN"
	JobStatusClosed JobStatus = "CLOSED"
)

type Job struct {
	ID                 uuid.UUID      `gorm:"column:id;type:uuid;primaryKey"`
	CompanyId          uuid.UUID      `gorm:"column:company_id;type:uuid;not null"`
	Title              string         `gorm:"column:title;type:varchar;size:255;not null"`
	Description        string         `gorm:"column:description;type:text;not null"`
	Requirements       *string        `gorm:"column:requirements;type:text"`
	EmploymentType     EmploymentType `gorm:"column:employment_type;type:varchar;size:50;not null"`
	WorkMode           WorkMode       `gorm:"column:work_mode;type:varchar;size:50;not null"`
	Location           string         `gorm:"column:location;type:varchar;size:255;not null"`
	MinSalary          *int64         `gorm:"column:min_salary;type:bigint"`
	MaxSalary          *int64         `gorm:"column:max_salary;type:bigint"`
	Currency           string         `gorm:"column:currency;type:varchar;size:10;default:'IDR';not null"`
	IsSalaryNegotiable bool           `gorm:"column:is_salary_negotiable;type:boolean;default:false;not null"`
	Status             JobStatus      `gorm:"column:status;type:varchar;size:20;default:'OPEN';not null"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	Company            *Company       `gorm:"foreignKey:CompanyId;references:ID"`
}

func (j *Job) TableName() string {
	return "jobs"
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		j.ID = uuidV7
	}
	if j.Currency == "" {
		j.Currency = "IDR"
	}
	if j.Status == "" {
		j.Status = JobStatusOpen
	}
	return nil
}
