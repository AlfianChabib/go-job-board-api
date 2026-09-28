package repository_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCompanyRepository(t *testing.T) {
	ctx := context.Background()
	compID := uuid.New()
	recruiterID := uuid.New()

	t.Run("Create success and error", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCompanyRepository(db)

		mock.Expect(mockExpectation{
			match: "insert into \"companies\"",
			err:   errors.New("db error"),
		})
		res, err := repo.Create(ctx, domain.Company{Name: "Acme", RecruiterId: recruiterID})
		assert.Nil(t, res)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:        "insert into \"companies\"",
			rowsAffected: 1,
		})
		res, err = repo.Create(ctx, domain.Company{Name: "Acme", RecruiterId: recruiterID})
		require.NoError(t, err)
		assert.Equal(t, "Acme", res.Name)
	})

	t.Run("FindByRecruiterId, GetCompanyIdByRecruiterId, and FindById", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCompanyRepository(db)

		// FindByRecruiterId error & success
		mock.Expect(mockExpectation{match: "select * from \"companies\"", err: gorm.ErrRecordNotFound})
		_, err := repo.FindByRecruiterId(ctx, recruiterID)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "recruiter_id", "name"},
			rows:    [][]driver.Value{{compID.String(), recruiterID.String(), "Acme"}},
		})
		res, err := repo.FindByRecruiterId(ctx, recruiterID)
		require.NoError(t, err)
		assert.Equal(t, compID, res.ID)

		// GetCompanyIdByRecruiterId error & success
		mock.Expect(mockExpectation{match: "select \"id\" from \"companies\"", err: gorm.ErrRecordNotFound})
		id, err := repo.GetCompanyIdByRecruiterId(ctx, recruiterID)
		assert.Equal(t, uuid.Nil, id)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select \"id\" from \"companies\"",
			columns: []string{"id"},
			rows:    [][]driver.Value{{compID.String()}},
		})
		id, err = repo.GetCompanyIdByRecruiterId(ctx, recruiterID)
		require.NoError(t, err)
		assert.Equal(t, compID, id)

		// FindById error & success
		mock.Expect(mockExpectation{match: "select * from \"companies\"", err: gorm.ErrRecordNotFound})
		_, err = repo.FindById(ctx, compID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		res, err = repo.FindById(ctx, compID)
		require.NoError(t, err)
		assert.Equal(t, compID, res.ID)
	})

	t.Run("FindAll count error, empty total, find error, and success with filters", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCompanyRepository(db)
		req := web.GetCompaniesRequest{Search: "acme", Industry: "IT", Page: 0, Limit: 0}

		// 1. Count error
		mock.Expect(mockExpectation{match: "select count(*)", err: errors.New("count err")})
		_, _, err := repo.FindAll(ctx, req)
		assert.Error(t, err)

		// 2. Total == 0 early return
		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(0)}},
		})
		list, total, err := repo.FindAll(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, list)
		assert.Equal(t, int64(0), total)

		// 3. Find error
		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match: "select * from \"companies\"",
			err:   errors.New("find err"),
		})
		_, _, err = repo.FindAll(ctx, req)
		assert.Error(t, err)

		// 4. Success
		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		list, total, err = repo.FindAll(ctx, req)
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("Update, UpdateLogo, and UpdateBanner", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewCompanyRepository(db)

		// Update by ID error
		mock.Expect(mockExpectation{match: "update \"companies\"", err: errors.New("update err")})
		_, err := repo.Update(ctx, domain.Company{ID: compID, Name: "Acme"})
		assert.Error(t, err)

		// Update by RecruiterId 0 rows affected
		mock.Expect(mockExpectation{match: "update \"companies\"", rows: nil, rowsAffected: 0})
		_, err = repo.Update(ctx, domain.Company{RecruiterId: recruiterID, Name: "Acme"})
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Update success
		mock.Expect(mockExpectation{
			match:        "update \"companies\"",
			columns:      []string{"id", "name"},
			rows:         [][]driver.Value{{compID.String(), "Acme Updated"}},
			rowsAffected: 1,
		})
		updated, err := repo.Update(ctx, domain.Company{ID: compID, Name: "Acme Updated"})
		require.NoError(t, err)
		assert.Equal(t, "Acme Updated", updated.Name)

		// UpdateLogo error, 0 rows, and success
		mock.Expect(mockExpectation{match: "update \"companies\"", err: errors.New("db err")})
		assert.Error(t, repo.UpdateLogo(ctx, recruiterID, "logo.png"))

		mock.Expect(mockExpectation{match: "update \"companies\"", rowsAffected: 0})
		assert.ErrorIs(t, repo.UpdateLogo(ctx, recruiterID, "logo.png"), gorm.ErrRecordNotFound)

		mock.Expect(mockExpectation{match: "update \"companies\"", rowsAffected: 1})
		assert.NoError(t, repo.UpdateLogo(ctx, recruiterID, "logo.png"))

		// UpdateBanner error, 0 rows, and success
		mock.Expect(mockExpectation{match: "update \"companies\"", err: errors.New("db err")})
		assert.Error(t, repo.UpdateBanner(ctx, recruiterID, "banner.png"))

		mock.Expect(mockExpectation{match: "update \"companies\"", rowsAffected: 0})
		assert.ErrorIs(t, repo.UpdateBanner(ctx, recruiterID, "banner.png"), gorm.ErrRecordNotFound)

		mock.Expect(mockExpectation{match: "update \"companies\"", rowsAffected: 1})
		assert.NoError(t, repo.UpdateBanner(ctx, recruiterID, "banner.png"))
	})
}
