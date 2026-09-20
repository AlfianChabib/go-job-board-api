package web

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"io"

	"github.com/google/uuid"
)

type CreateCompanyRequest struct {
	Name         string              `json:"name" validate:"required,max=255"`
	Location     string              `json:"location" validate:"required,max=255"`
	Industry     string              `json:"industry" validate:"required,max=255"`
	EmployeeSize domain.EmployeeSize `json:"employee_size" validate:"required,oneof=1-50 51-200 201-500 500+"`
	Description  *string             `json:"description,omitempty"`
	Website      *string             `json:"website,omitempty" validate:"omitempty,nullable_url"`
}

type UpdateCompanyRequest struct {
	Name         string              `json:"name" validate:"required,max=255"`
	Location     string              `json:"location" validate:"required,max=255"`
	Industry     string              `json:"industry" validate:"required,max=255"`
	EmployeeSize domain.EmployeeSize `json:"employee_size" validate:"required,oneof=1-50 51-200 201-500 500+"`
	Description  *string             `json:"description,omitempty"`
	Website      *string             `json:"website,omitempty" validate:"omitempty,nullable_url"`
}

type UpdateCompanyLogoRequest struct {
	RecruiterId uuid.UUID
	File        io.Reader
	FileSize    int64
	ContentType string
	Extension   string
}

type UpdateCompanyBannerRequest struct {
	RecruiterId uuid.UUID
	File        io.Reader
	FileSize    int64
	ContentType string
	Extension   string
}

type GetCompaniesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search" validate:"omitempty"`
	Industry string `query:"industry" validate:"omitempty"`
}

type GetCompanyByIdRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}
