package web

import (
	"io"

	"github.com/google/uuid"
)

type CandidateRequest struct {
	UserId uuid.UUID `json:"user_id"`
}

type UpdateCandidateRequest struct {
	Headline string `json:"headline" validate:"required,max=255"`
	Phone    string `json:"phone" validate:"required,e164"`
}

type UpdateCandidateAvatarRequest struct {
	UserId     uuid.UUID
	File       io.Reader
	FileSize   int64
	ContenType string
	Extension  string
}

type UpdateCandidateSkillsRequest struct {
	Skills []SkillRequest `json:"skills" validate:"required,min=1,dive"`
}

type SkillRequest struct {
	Id   *uuid.UUID `json:"id,omitempty"`
	Name string     `json:"name" validate:"required"`
}

type CreateExperienceRequest struct {
	CompanyName string  `json:"company_name" validate:"required"`
	Position    string  `json:"position" validate:"required"`
	StartDate   string  `json:"start_date" validate:"required,rfc3339"`
	EndDate     *string `json:"end_date" validate:"required_if=IsCurrent false,rfc3339"`
	IsCurrent   bool    `json:"is_current"`
	Description *string `json:"description,omitempty"`
}

type UpdateExperienceRequest struct {
	ExperienceId uuid.UUID `param:"id" validate:"required"`
	CompanyName  string    `json:"company_name" validate:"required"`
	Position     string    `json:"position" validate:"required"`
	StartDate    string    `json:"start_date" validate:"required,rfc3339"`
	EndDate      *string   `json:"end_date" validate:"required_if=IsCurrent false,rfc3339"`
	IsCurrent    bool      `json:"is_current"`
	Description  *string   `json:"description,omitempty"`
}
