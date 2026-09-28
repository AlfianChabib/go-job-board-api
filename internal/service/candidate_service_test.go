package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCandidateService_Get(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("Not found returns Candidate not found error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("Get", ctx, userID).Return(nil, errors.New("db error"))

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.Get(ctx, userID)
		assert.Nil(t, res)
		assert.EqualError(t, err, "Candidate not found")
	})

	t.Run("Success returns GetCandidateResponse", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		profileID := uuid.New()
		headline := "Backend Developer"
		repo.On("Get", ctx, userID).Return(&domain.Profile{
			ID:       profileID,
			UserId:   userID,
			Headline: &headline,
		}, nil)

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.Get(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profileID, res.Id)
		assert.Equal(t, userID, res.UserId)
		assert.Equal(t, &headline, res.Headline)
	})
}

func TestCandidateService_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	headline := "Senior Go Developer"
	phone := "+628123456789"
	input := domain.Profile{
		UserId:   userID,
		Headline: &headline,
		Phone:    &phone,
	}

	t.Run("Repository error returns error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("Update", ctx, input).Return(nil, errors.New("update err"))

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.Update(ctx, input)
		assert.Nil(t, res)
		assert.EqualError(t, err, "update err")
	})

	t.Run("Success returns UpdateCandidateResponse", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		profileID := uuid.New()
		repo.On("Update", ctx, input).Return(&domain.Profile{
			ID:       profileID,
			UserId:   userID,
			Headline: &headline,
			Phone:    &phone,
		}, nil)

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.Update(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, profileID, res.Id)
		assert.Equal(t, userID, res.UserId)
		assert.Equal(t, &headline, res.Headline)
		assert.Equal(t, &phone, res.Phone)
	})
}

func TestCandidateService_UploadAvatar(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	req := web.UpdateCandidateAvatarRequest{
		UserId:     userID,
		File:       bytes.NewReader([]byte("img")),
		FileSize:   3,
		ContenType: "image/png",
		Extension:  ".png",
	}

	t.Run("Storage UploadAvatar error returns error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		storage.On("UploadAvatar", ctx, mock.AnythingOfType("string"), req.File, req.FileSize, req.ContenType).
			Return(nil, errors.New("minio error"))

		svc := NewCandidateService(repo, storage)
		res, err := svc.UploadAvatar(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "minio error")
	})

	t.Run("Repository UploadAvatar error returns error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		url := "http://localhost:9000/avatar/file.png"
		storage.On("UploadAvatar", ctx, mock.AnythingOfType("string"), req.File, req.FileSize, req.ContenType).
			Return(&url, nil)
		repo.On("UploadAvatar", ctx, userID, url).Return(errors.New("db error"))

		svc := NewCandidateService(repo, storage)
		res, err := svc.UploadAvatar(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Success returns avatar URL", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		url := "http://localhost:9000/avatar/file.png"
		storage.On("UploadAvatar", ctx, mock.AnythingOfType("string"), req.File, req.FileSize, req.ContenType).
			Return(&url, nil)
		repo.On("UploadAvatar", ctx, userID, url).Return(nil)

		svc := NewCandidateService(repo, storage)
		res, err := svc.UploadAvatar(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, &url, res)
	})
}

func TestCandidateService_DeleteAvatar(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	avatarUrl := "http://localhost:9000/avatar/my_avatar.png"

	t.Run("Repository Get error returns internal server error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("Get", ctx, userID).Return(nil, errors.New("not found"))

		svc := NewCandidateService(repo, new(mockStorageRepository))
		err := svc.DeleteAvatar(ctx, userID)
		assert.EqualError(t, err, "internal server error")
	})

	t.Run("Storage DeleteAvatar error returns internal server error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		repo.On("Get", ctx, userID).Return(&domain.Profile{AvatarUrl: &avatarUrl}, nil)
		storage.On("DeleteAvatar", ctx, "my_avatar.png").Return(errors.New("minio delete err"))

		svc := NewCandidateService(repo, storage)
		err := svc.DeleteAvatar(ctx, userID)
		assert.EqualError(t, err, "internal server error")
	})

	t.Run("Repository DeleteAvatar error returns internal server error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		repo.On("Get", ctx, userID).Return(&domain.Profile{AvatarUrl: &avatarUrl}, nil)
		storage.On("DeleteAvatar", ctx, "my_avatar.png").Return(nil)
		repo.On("DeleteAvatar", ctx, userID).Return(errors.New("db delete err"))

		svc := NewCandidateService(repo, storage)
		err := svc.DeleteAvatar(ctx, userID)
		assert.EqualError(t, err, "internal server error")
	})

	t.Run("Success deletes avatar from storage and db", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		storage := new(mockStorageRepository)
		repo.On("Get", ctx, userID).Return(&domain.Profile{AvatarUrl: &avatarUrl}, nil)
		storage.On("DeleteAvatar", ctx, "my_avatar.png").Return(nil)
		repo.On("DeleteAvatar", ctx, userID).Return(nil)

		svc := NewCandidateService(repo, storage)
		err := svc.DeleteAvatar(ctx, userID)
		assert.NoError(t, err)
	})
}

func TestCandidateService_UpdateSkills(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	req := web.UpdateCandidateSkillsRequest{
		Skills: []web.SkillRequest{{Name: "Go"}},
	}

	t.Run("Record not found returns 404 fiber error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("UpdateSkills", ctx, userID, req).Return(nil, gorm.ErrRecordNotFound)

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.UpdateSkills(ctx, userID, req)
		assert.Nil(t, res)
		var fiberErr *fiber.Error
		require.True(t, errors.As(err, &fiberErr))
		assert.Equal(t, fiber.StatusNotFound, fiberErr.Code)
	})

	t.Run("Generic error returns internal server error", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("UpdateSkills", ctx, userID, req).Return(nil, errors.New("some db err"))

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.UpdateSkills(ctx, userID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "internal server error")
	})

	t.Run("Success returns updated skills", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		skills := []domain.Skill{{ID: uuid.New(), Name: "go"}}
		repo.On("UpdateSkills", ctx, userID, req).Return(skills, nil)

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.UpdateSkills(ctx, userID, req)
		require.NoError(t, err)
		assert.Equal(t, skills, res)
	})
}

func TestCandidateService_Experiences(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	profileID := uuid.New()
	endDateStr := "2026-01-01T00:00:00Z"

	t.Run("GetExperiences error and success", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		repo.On("GetExperiences", ctx, userID).Return(nil, errors.New("db err")).Once()

		svc := NewCandidateService(repo, new(mockStorageRepository))
		res, err := svc.GetExperiences(ctx, userID)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db err")

		exps := []domain.Experience{{ID: uuid.New(), CompanyName: "Acme"}}
		repo.On("GetExperiences", ctx, userID).Return(exps, nil).Once()
		res, err = svc.GetExperiences(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, exps, res)
	})

	t.Run("CreateExperience error and success", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		req := web.CreateExperienceRequest{
			CompanyName: "Acme",
			Position:    "Developer",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     &endDateStr,
		}

		repo.On("Get", ctx, userID).Return(&domain.Profile{ID: profileID}, nil)
		repo.On("CreateExperience", ctx, userID, mock.AnythingOfType("domain.Experience")).
			Return(errors.New("create err")).Once()

		svc := NewCandidateService(repo, new(mockStorageRepository))
		err := svc.CreateExperience(ctx, userID, req)
		assert.EqualError(t, err, "create err")

		repo.On("CreateExperience", ctx, userID, mock.AnythingOfType("domain.Experience")).
			Return(nil).Once()
		err = svc.CreateExperience(ctx, userID, req)
		assert.NoError(t, err)
	})

	t.Run("UpdateExperience profile error and success", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		expID := uuid.New()
		req := web.UpdateExperienceRequest{
			ExperienceId: expID,
			CompanyName:  "Acme",
			Position:     "Lead",
			StartDate:    "2025-01-01T00:00:00Z",
			EndDate:      &endDateStr,
		}

		repo.On("GetProfileIdByUserId", ctx, userID).Return(uuid.Nil, errors.New("profile not found")).Once()
		svc := NewCandidateService(repo, new(mockStorageRepository))
		err := svc.UpdateExperience(ctx, userID, req)
		assert.EqualError(t, err, "profile not found")

		repo.On("GetProfileIdByUserId", ctx, userID).Return(profileID, nil).Once()
		repo.On("UpdateExperience", ctx, expID, mock.AnythingOfType("domain.Experience")).Return(nil).Once()
		err = svc.UpdateExperience(ctx, userID, req)
		assert.NoError(t, err)
	})

	t.Run("DeleteExperience profile error and success", func(t *testing.T) {
		repo := new(mockCandidateRepository)
		expID := uuid.New()

		repo.On("GetProfileIdByUserId", ctx, userID).Return(uuid.Nil, errors.New("profile not found")).Once()
		svc := NewCandidateService(repo, new(mockStorageRepository))
		err := svc.DeleteExperience(ctx, userID, expID)
		assert.EqualError(t, err, "profile not found")

		repo.On("GetProfileIdByUserId", ctx, userID).Return(profileID, nil).Once()
		repo.On("DeleteExperience", ctx, profileID, expID).Return(nil).Once()
		err = svc.DeleteExperience(ctx, userID, expID)
		assert.NoError(t, err)
	})
}
