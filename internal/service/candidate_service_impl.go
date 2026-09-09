package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type candidateService struct {
	repository repository.CandidateRepository
}

func NewCandidateService(repo repository.CandidateRepository) CandidateService {
	return &candidateService{
		repository: repo,
	}
}

func (service *candidateService) Get(ctx context.Context, userId uuid.UUID) (*web.GetCandidateResponse, error) {
	candidate, err := service.repository.Get(ctx, userId)
	if err != nil {
		return nil, errors.New("Candidate not found")
	}
	return &web.GetCandidateResponse{
		Id:          candidate.ID,
		UserId:      candidate.UserId,
		AvatarUrl:   candidate.AvatarUrl,
		Headline:    candidate.Headline,
		Bio:         candidate.Bio,
		Phone:       candidate.Phone,
		ResumeUrl:   candidate.ResumeUrl,
		CreatedAt:   candidate.CreatedAt,
		UpdatedAt:   candidate.UpdatedAt,
		Skills:      candidate.Skills,
		Experiences: candidate.Experiences,
	}, nil
}

func (service *candidateService) Update(ctx context.Context, candidate domain.Profile) (*web.UpdateCandidateResponse, error) {
	data, err := service.repository.Update(ctx, candidate)
	if err != nil {
		return nil, err
	}
	return &web.UpdateCandidateResponse{
		Id:       data.ID,
		UserId:   data.UserId,
		Headline: data.Headline,
		Phone:    data.Phone,
	}, nil
}
