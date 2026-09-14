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

		var skillsWithId []domain.Skill
		var newSkillNames []string
		for _, req := range skills.Skills {
			if req.Id != nil && *req.Id != uuid.Nil {
				skillsWithId = append(skillsWithId, domain.Skill{
					ID:   *req.Id,
					Name: req.Name,
				})
			} else {
				newSkillNames = append(newSkillNames, strings.ToLower(req.Name))
			}
		}

		if len(skillsWithId) == 0 && len(newSkillNames) == 0 {
			if err := tx.Model(&profile).Association("Skills").Clear(); err != nil {
				return err
			}
			finalSkills = []domain.Skill{}
			return nil
		}

		if len(newSkillNames) > 0 {
			var newSkillsToInsert []domain.Skill
			for _, name := range newSkillNames {
				newSkillsToInsert = append(newSkillsToInsert, domain.Skill{Name: name, Label: ""})
			}

			if err := tx.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "name"}},
					DoNothing: true,
				},
			).Create(newSkillsToInsert).Error; err != nil {
				return err
			}

			var resolvedNewSkills []domain.Skill
			if err := tx.Where("name IN ?", newSkillNames).Find(&resolvedNewSkills).Error; err != nil {
				return err
			}

			skillsWithId = append(skillsWithId, resolvedNewSkills...)
		}

		skillsId := make([]uuid.UUID, 0, len(skillsWithId))
		for _, skill := range skillsWithId {
			skillsId = append(skillsId, skill.ID)
		}

		if err := tx.Where("id IN ?", skillsId).Find(&finalSkills).Error; err != nil {
			return err
		}

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
