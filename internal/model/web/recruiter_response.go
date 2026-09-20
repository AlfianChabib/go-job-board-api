package web

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"time"

	"github.com/google/uuid"
)

// CompanyResponse represents common company profile data
// Digunakan pada:
// - POST /api/recruiter/company
// - GET  /api/recruiter/company
// - PUT  /api/recruiter/company
type CompanyResponse struct {
	ID           uuid.UUID           `json:"id"`
	RecruiterId  uuid.UUID           `json:"recruiter_id"`
	Name         string              `json:"name"`
	LogoUrl      *string             `json:"logo_url"`
	BannerUrl    *string             `json:"banner_url"`
	Website      *string             `json:"website"`
	Industry     string              `json:"industry"`
	EmployeeSize domain.EmployeeSize `json:"employee_size"`
	Description  *string             `json:"description"`
	Location     string              `json:"location"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// UploadCompanyLogoResponse response untuk PUT /api/recruiter/company/logo
type UploadCompanyLogoResponse struct {
	LogoUrl string `json:"logo_url"`
}

// UploadCompanyBannerResponse response untuk PUT /api/recruiter/company/banner
type UploadCompanyBannerResponse struct {
	BannerUrl string `json:"banner_url"`
}

// CompanyPaginationResponse response untuk GET /api/companies (Public/Pagination)
type CompanyPaginationResponse struct {
	Companies  []CompanyResponse `json:"companies"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// CompanyDetailResponse response untuk GET /api/companies/:id (Detail profil beserta lowongannya)
type CompanyDetailResponse struct {
	CompanyResponse
	Jobs []CompanyJobResponse `json:"jobs"`
}

// CompanyJobResponse item ringkasan lowongan pekerjaan pada detail company
type CompanyJobResponse struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Requirements       *string   `json:"requirements"`
	EmploymentType     string    `json:"employment_type"`
	WorkMode           string    `json:"work_mode"`
	Location           string    `json:"location"`
	MinSalary          *int64    `json:"min_salary"`
	MaxSalary          *int64    `json:"max_salary"`
	Currency           string    `json:"currency"`
	IsSalaryNegotiable bool      `json:"is_salary_negotiable"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
