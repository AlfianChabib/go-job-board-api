package web

import "github.com/google/uuid"

type CreateJobRequest struct {
	Title              string  `json:"title" validate:"required,min=3,max=255"`
	Description        string  `json:"description" validate:"required,min=10"`
	Requirements       *string `json:"requirements,omitempty" validate:"omitempty"`
	EmploymentType     string  `json:"employment_type" validate:"required,oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE"`
	WorkMode           string  `json:"work_mode" validate:"required,oneof=ON_SITE HYBRID REMOTE"`
	Location           string  `json:"location" validate:"required,max=255"`
	MinSalary          *int64  `json:"min_salary,omitempty" validate:"omitempty,min=0"`
	MaxSalary          *int64  `json:"max_salary,omitempty" validate:"omitempty,gtefield=MinSalary"`
	Currency           string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	IsSalaryNegotiable bool    `json:"is_salary_negotiable"`
}

type UpdateJobRequest struct {
	ID                 uuid.UUID `param:"id" validate:"required"`
	Title              string    `json:"title" validate:"required,min=3,max=255"`
	Description        string    `json:"description" validate:"required,min=10"`
	Requirements       *string   `json:"requirements,omitempty" validate:"omitempty"`
	EmploymentType     string    `json:"employment_type" validate:"required,oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE"`
	WorkMode           string    `json:"work_mode" validate:"required,oneof=ON_SITE HYBRID REMOTE"`
	Location           string    `json:"location" validate:"required,max=255"`
	MinSalary          *int64    `json:"min_salary,omitempty" validate:"omitempty,min=0"`
	MaxSalary          *int64    `json:"max_salary,omitempty" validate:"omitempty,gtefield=MinSalary"`
	Currency           string    `json:"currency,omitempty" validate:"omitempty,len=3"`
	IsSalaryNegotiable bool      `json:"is_salary_negotiable"`
}

type UpdateJobStatusRequest struct {
	ID     uuid.UUID `param:"id" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=OPEN CLOSED"`
}

type GetJobsRequest struct {
	Page           int    `query:"page" validate:"omitempty,min=1"`
	Limit          int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search         string `query:"search" validate:"omitempty"`
	EmploymentType string `query:"employment_type" validate:"omitempty,oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE"`
	WorkMode       string `query:"work_mode" validate:"omitempty,oneof=ON_SITE HYBRID REMOTE"`
	Location       string `query:"location" validate:"omitempty"`
	MinSalary      *int64 `query:"min_salary" validate:"omitempty,min=0"`
}

type GetRecruiterJobsRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Status string `query:"status" validate:"omitempty,oneof=OPEN CLOSED"`
}

type GetJobByIdRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

type DeleteJobRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}
