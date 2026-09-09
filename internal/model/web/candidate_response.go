package web

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"time"

	"github.com/google/uuid"
)

type GetCandidateResponse struct {
	Id          uuid.UUID           `json:"id"`
	UserId      uuid.UUID           `json:"user_id"`
	AvatarUrl   *string             `json:"avatar_url"`
	Headline    *string             `json:"headline"`
	Bio         *string             `json:"bio"`
	Phone       *string             `json:"phone"`
	ResumeUrl   *string             `json:"resume_url"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Skills      []domain.Skill      `json:"skills"`
	Experiences []domain.Experience `json:"experiences"`
}

type UpdateCandidateResponse struct {
	Id       uuid.UUID `json:"id"`
	UserId   uuid.UUID `json:"user_id"`
	Headline *string   `json:"headline"`
	Phone    *string   `json:"phone"`
}
