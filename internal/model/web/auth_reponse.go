package web

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"

	"github.com/google/uuid"
)

type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type RefreshTokenResponse struct {
	*domain.TokenPair
}
