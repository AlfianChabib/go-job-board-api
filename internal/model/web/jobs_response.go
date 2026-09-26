package web

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"time"

	"github.com/google/uuid"
)

type JobCompanySummary struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	LogoUrl  *string   `json:"logo_url"`
	Location string    `json:"location"`
}

type JobResponse struct {
	ID                 uuid.UUID             `json:"id"`
	CompanyId          uuid.UUID             `json:"company_id"`
	Title              string                `json:"title"`
	Description        string                `json:"description"`
	Requirements       *string               `json:"requirements"`
	EmploymentType     domain.EmploymentType `json:"employment_type"`
	WorkMode           domain.WorkMode       `json:"work_mode"`
	Location           string                `json:"location"`
	MinSalary          *int64                `json:"min_salary"`
	MaxSalary          *int64                `json:"max_salary"`
	Currency           string                `json:"currency"`
	IsSalaryNegotiable bool                  `json:"is_salary_negotiable"`
	Status             domain.JobStatus      `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	Company            *JobCompanySummary    `json:"company,omitempty"`
}

type JobPaginationResponse struct {
	Jobs       []JobResponse `json:"jobs"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
