package repository_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCandidateRepository(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	profileID := uuid.New()
	expID := uuid.New()

	t.Run("Get and GetProfileIdByUserId", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCandidateRepository(db)

		// Get error & success
		mock.Expect(mockExpectation{match: "select * from \"profiles\"", err: errors.New("not found")})
		_, err := repo.Get(ctx, userID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), userID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"experiences\"",
			columns: []string{"id", "profile_id"},
			rows:    nil,
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"profile_skills\"",
			columns: []string{"profile_id", "skill_id"},
			rows:    nil,
		})
		prof, err := repo.Get(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profileID, prof.ID)

		// GetProfileIdByUserId error & success
		mock.Expect(mockExpectation{match: "select \"id\" from \"profiles\"", err: errors.New("not found")})
		id, err := repo.GetProfileIdByUserId(ctx, userID)
		assert.Equal(t, uuid.Nil, id)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select \"id\" from \"profiles\"",
			columns: []string{"id"},
			rows:    [][]driver.Value{{profileID.String()}},
		})
		id, err = repo.GetProfileIdByUserId(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profileID, id)
	})

	t.Run("Update, UploadAvatar, and DeleteAvatar", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCandidateRepository(db)
		headline := "Go Dev"

		// Update error & success
		mock.Expect(mockExpectation{match: "update \"profiles\"", err: errors.New("update err")})
		_, err := repo.Update(ctx, domain.Profile{UserId: userID, Headline: &headline})
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:        "update \"profiles\"",
			columns:      []string{"id", "user_id", "headline"},
			rows:         [][]driver.Value{{profileID.String(), userID.String(), headline}},
			rowsAffected: 1,
		})
		res, err := repo.Update(ctx, domain.Profile{UserId: userID, Headline: &headline})
		require.NoError(t, err)
		assert.Equal(t, &headline, res.Headline)

		// UploadAvatar error & success
		mock.Expect(mockExpectation{match: "update \"profiles\"", err: errors.New("db err")})
		assert.Error(t, repo.UploadAvatar(ctx, userID, "http://avatar.png"))

		mock.Expect(mockExpectation{match: "update \"profiles\"", rowsAffected: 1})
		assert.NoError(t, repo.UploadAvatar(ctx, userID, "http://avatar.png"))

		// DeleteAvatar error & success
		mock.Expect(mockExpectation{match: "update \"profiles\"", err: errors.New("db err")})
		assert.Error(t, repo.DeleteAvatar(ctx, userID))

		mock.Expect(mockExpectation{match: "update \"profiles\"", rowsAffected: 1})
		assert.NoError(t, repo.DeleteAvatar(ctx, userID))
	})

	t.Run("UpdateSkills scenarios", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCandidateRepository(db)

		// 1. Profile not found
		mock.Expect(mockExpectation{match: "select * from \"profiles\"", err: errors.New("not found")})
		_, err := repo.UpdateSkills(ctx, userID, web.UpdateCandidateSkillsRequest{})
		assert.Error(t, err)

		// 2. Empty skills clears association
		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), userID.String()}},
		})
		mock.Expect(mockExpectation{match: "delete from \"profile_skills\"", rowsAffected: 1})
		skills, err := repo.UpdateSkills(ctx, userID, web.UpdateCandidateSkillsRequest{Skills: nil})
		require.NoError(t, err)
		assert.Empty(t, skills)

		// 3. Existing skill with ID + new skill without ID
		existingSkillID := uuid.New()
		newSkillID := uuid.New()
		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), userID.String()}},
		})
		mock.Expect(mockExpectation{match: "insert into \"skills\"", rowsAffected: 1})
		mock.Expect(mockExpectation{
			match:   "select * from \"skills\" where name in",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{newSkillID.String(), "docker"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"skills\" where id in",
			columns: []string{"id", "name"},
			rows: [][]driver.Value{
				{existingSkillID.String(), "go"},
				{newSkillID.String(), "docker"},
			},
		})
		mock.Expect(mockExpectation{match: "update \"profiles\"", rowsAffected: 1})
		mock.Expect(mockExpectation{match: "insert into \"skills\"", rowsAffected: 2})
		mock.Expect(mockExpectation{match: "insert into \"profile_skills\"", rowsAffected: 2})
		mock.Expect(mockExpectation{match: "delete from \"profile_skills\"", rowsAffected: 0})

		skills, err = repo.UpdateSkills(ctx, userID, web.UpdateCandidateSkillsRequest{
			Skills: []web.SkillRequest{
				{Id: &existingSkillID, Name: "go"},
				{Name: "Docker"},
			},
		})
		require.NoError(t, err)
		assert.Len(t, skills, 2)
	})

	t.Run("GetExperiences, CreateExperience, UpdateExperience, and DeleteExperience", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCandidateRepository(db)

		// GetExperiences error & success
		mock.Expect(mockExpectation{match: "select * from \"profiles\"", err: errors.New("db err")})
		_, err := repo.GetExperiences(ctx, userID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), userID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"experiences\"",
			columns: []string{"id", "profile_id", "company_name"},
			rows:    [][]driver.Value{{expID.String(), profileID.String(), "Acme"}},
		})
		exps, err := repo.GetExperiences(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, exps, 1)

		// CreateExperience error & success
		mock.Expect(mockExpectation{match: "insert into \"experiences\"", err: errors.New("db err")})
		assert.Error(t, repo.CreateExperience(ctx, userID, domain.Experience{ProfileId: profileID, CompanyName: "Acme"}))

		mock.Expect(mockExpectation{match: "insert into \"experiences\"", rowsAffected: 1})
		assert.NoError(t, repo.CreateExperience(ctx, userID, domain.Experience{ProfileId: profileID, CompanyName: "Acme"}))

		// UpdateExperience error, 0 rows, and success
		mock.Expect(mockExpectation{match: "update \"experiences\"", err: errors.New("db err")})
		assert.Error(t, repo.UpdateExperience(ctx, expID, domain.Experience{ProfileId: profileID, CompanyName: "Acme"}))

		mock.Expect(mockExpectation{match: "update \"experiences\"", rowsAffected: 0})
		assert.Equal(t, errs.ErrExperienceNotFound, repo.UpdateExperience(ctx, expID, domain.Experience{ProfileId: profileID, CompanyName: "Acme"}))

		mock.Expect(mockExpectation{match: "update \"experiences\"", rowsAffected: 1})
		assert.NoError(t, repo.UpdateExperience(ctx, expID, domain.Experience{ProfileId: profileID, CompanyName: "Acme"}))

		// DeleteExperience error, 0 rows, and success
		mock.Expect(mockExpectation{match: "delete from \"experiences\"", err: errors.New("db err")})
		assert.Error(t, repo.DeleteExperience(ctx, profileID, expID))

		mock.Expect(mockExpectation{match: "delete from \"experiences\"", rowsAffected: 0})
		assert.Equal(t, errs.ErrExperienceNotFound, repo.DeleteExperience(ctx, profileID, expID))

		mock.Expect(mockExpectation{match: "delete from \"experiences\"", rowsAffected: 1})
		assert.NoError(t, repo.DeleteExperience(ctx, profileID, expID))
	})
}
