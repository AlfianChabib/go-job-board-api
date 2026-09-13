package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type candidateRepository struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{
		db: db,
	}
}

func (repo *candidateRepository) Get(ctx context.Context, userId uuid.UUID) (*domain.Profile, error) {
	var candidate domain.Profile
	err := repo.db.WithContext(ctx).
		Preload("Skills").
		Preload("Experiences").
		Take(&candidate, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

func (repo *candidateRepository) Update(ctx context.Context, candidate domain.Profile) (*domain.Profile, error) {
	err := repo.db.WithContext(ctx).
		Model(&candidate).
		Clauses(clause.Returning{}). // <-- Isi otomatis ID & kolom lainnya dari DB
		Where("user_id = ?", candidate.UserId).
		Select("headline", "phone").
		Updates(&candidate).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (repo *candidateRepository) UploadAvatar(ctx context.Context, userId uuid.UUID, avatarUrl string) error {
	err := repo.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("user_id = ?", userId).
		Update("avatar_url", avatarUrl).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *candidateRepository) DeleteAvatar(ctx context.Context, userId uuid.UUID) error {
	err := repo.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("user_id = ?", userId).
		Update("avatar_url", nil).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *candidateRepository) UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) (*[]domain.Skill, error) {
	var finalSkills []domain.Skill

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var profile domain.Profile
		if err := tx.Where("user_id = ?", userId).Take(&profile).Error; err != nil {
			return err
		}

		var uniqueNames []string
		var lowerNames []string
		seen := make(map[string]bool)

		for _, s := range skills.Skills {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				continue
			}
			lower := strings.ToLower(trimmed)
			if !seen[lower] {
				seen[lower] = true
				uniqueNames = append(uniqueNames, trimmed)
				lowerNames = append(lowerNames, lower)
			}
		}

		if len(uniqueNames) == 0 {
			if err := tx.Model(&profile).Association("Skills").Clear(); err != nil {
				return err
			}
			finalSkills = []domain.Skill{}
			return nil
		}

		// Look for existing skills in database
		var existingSkills []domain.Skill
		if err := tx.Where("LOWER(name) IN ?", lowerNames).Find(&existingSkills).Error; err != nil {
			return err
		}

		existingMap := make(map[string]domain.Skill)
		for _, es := range existingSkills {
			existingMap[strings.ToLower(es.Name)] = es
		}

		// Insert new skills if they do not exist
		var newSkills []domain.Skill
		for _, name := range uniqueNames {
			lower := strings.ToLower(name)
			if _, exists := existingMap[lower]; !exists {
				newSkills = append(newSkills, domain.Skill{
					Name:  name,
					Label: name,
				})
			}
		}

		if len(newSkills) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoNothing: true,
			}).Create(&newSkills).Error; err != nil {
				return err
			}
		}

		// Retrieve all matching skills
		if err := tx.Where("LOWER(name) IN ?", lowerNames).Find(&finalSkills).Error; err != nil {
			return err
		}

		// Replace candidate's skills association
		if err := tx.Model(&profile).Association("Skills").Replace(finalSkills); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &finalSkills, nil
}
