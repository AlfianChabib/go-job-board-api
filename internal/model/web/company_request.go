package web

import "github.com/google/uuid"

type GetCompaniesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search" validate:"omitempty"`
	Industry string `query:"industry" validate:"omitempty"`
}

type GetCompanyByIdRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}
