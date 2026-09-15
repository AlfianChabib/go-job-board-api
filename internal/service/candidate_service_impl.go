package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type candidateService struct {
	repository  repository.CandidateRepository
	storageRepo repository.StorageRepository
}

func NewCandidateService(repo repository.CandidateRepository, storageRepo repository.StorageRepository) CandidateService {
	return &candidateService{
		repository:  repo,
		storageRepo: storageRepo,
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

func (service *candidateService) UploadAvatar(ctx context.Context, req web.UpdateCandidateAvatarRequest) (*string, error) {
	fileName := fmt.Sprintf("%s_%d%s", req.UserId, time.Now().Unix(), req.Extension)

	avatarUrl, err := service.storageRepo.UploadAvatar(ctx, fileName, req.File, req.FileSize, req.ContenType)
	if err != nil {
		return nil, err
	}

	err = service.repository.UploadAvatar(ctx, req.UserId, *avatarUrl)
	if err != nil {
		return nil, err
	}

	return avatarUrl, nil
}

func (service *candidateService) DeleteAvatar(ctx context.Context, userId uuid.UUID) error {
	candidate, err := service.repository.Get(ctx, userId)
	if err != nil {
		return errors.New("internal server error")
	}

	avatarUrl := *candidate.AvatarUrl
	avatarFileName := strings.Split(avatarUrl, "avatar/")[1]
	err = service.storageRepo.DeleteAvatar(ctx, avatarFileName)
	if err != nil {
		return errors.New("internal server error")
	}

	err = service.repository.DeleteAvatar(ctx, userId)
	if err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func (service *candidateService) UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) (*[]domain.Skill, error) {
	updatedSkills, err := service.repository.UpdateSkills(ctx, userId, skills)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Candidate profile not found")
		}
		return nil, errors.New("internal server error")
	}

	return updatedSkills, nil
}

func (service *candidateService) GetExperiences(ctx context.Context, userId uuid.UUID) (*[]domain.Experience, error) {
	experiences, err := service.repository.GetExperiences(ctx, userId)
	if err != nil {
		return nil, err
	}
	return experiences, nil
}
