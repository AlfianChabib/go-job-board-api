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
