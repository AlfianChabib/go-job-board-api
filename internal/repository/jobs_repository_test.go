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

func TestJobRepository(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New()
	compID := uuid.New()

	t.Run("Create error and success", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewJobRepository(db)

		mock.Expect(mockExpectation{match: "insert into \"jobs\"", err: errors.New("insert err")})
		_, err := repo.Create(ctx, domain.Job{Title: "Go Dev", CompanyId: compID})
		assert.Error(t, err)

		mock.Expect(mockExpectation{match: "insert into \"jobs\"", rowsAffected: 1})
		res, err := repo.Create(ctx, domain.Job{Title: "Go Dev", CompanyId: compID})
		require.NoError(t, err)
		assert.Equal(t, "Go Dev", res.Title)
	})

	t.Run("FindAll with all filters, count error, find error, and success", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewJobRepository(db)
		minSal := int64(5000000)
		req := web.GetJobsRequest{
			Search:         "go",
			EmploymentType: "FULL_TIME",
			WorkMode:       "REMOTE",
			Location:       "Jakarta",
			MinSalary:      &minSal,
			Page:           0,
			Limit:          0,
		}

		// 1. Count error
		mock.Expect(mockExpectation{match: "select count(*)", err: errors.New("count err")})
		_, _, err := repo.FindAll(ctx, req)
		assert.Error(t, err)

		// 2. Find error
		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", err: errors.New("find err")})
		_, _, err = repo.FindAll(ctx, req)
		assert.Error(t, err)

		// 3. Success with Preload("Company")
		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"jobs\"",
			columns: []string{"id", "company_id", "title"},
			rows:    [][]driver.Value{{jobID.String(), compID.String(), "Go Dev"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		jobs, total, err := repo.FindAll(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, jobs, 1)
		assert.Equal(t, jobID, jobs[0].ID)
	})

	t.Run("FindAllByCompanyId count error, find error, and success", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewJobRepository(db)
		req := web.GetRecruiterJobsRequest{Status: "OPEN"}

		mock.Expect(mockExpectation{match: "select count(*)", err: errors.New("count err")})
		_, _, err := repo.FindAllByCompanyId(ctx, compID, req)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", err: errors.New("find err")})
		_, _, err = repo.FindAllByCompanyId(ctx, compID, req)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"jobs\"",
			columns: []string{"id", "company_id", "title"},
			rows:    [][]driver.Value{{jobID.String(), compID.String(), "Go Dev"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		jobs, total, err := repo.FindAllByCompanyId(ctx, compID, req)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, jobs, 1)
	})

	t.Run("FindById, Update, UpdateStatus, and Delete", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewJobRepository(db)

		// FindById error & success
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", err: gorm.ErrRecordNotFound})
		_, err := repo.FindById(ctx, jobID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"jobs\"",
			columns: []string{"id", "company_id", "title"},
			rows:    [][]driver.Value{{jobID.String(), compID.String(), "Go Dev"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		job, err := repo.FindById(ctx, jobID)
		require.NoError(t, err)
		assert.Equal(t, jobID, job.ID)

		// Update error & success
		mock.Expect(mockExpectation{match: "update \"jobs\"", err: errors.New("update err")})
		_, err = repo.Update(ctx, domain.Job{ID: jobID, Title: "New"})
		assert.Error(t, err)

		mock.Expect(mockExpectation{match: "update \"jobs\"", rowsAffected: 1})
		updated, err := repo.Update(ctx, domain.Job{ID: jobID, Title: "New"})
		require.NoError(t, err)
		assert.Equal(t, "New", updated.Title)

		// UpdateStatus error, 0 rows, and success
		mock.Expect(mockExpectation{match: "update \"jobs\"", err: errors.New("db err")})
		assert.Error(t, repo.UpdateStatus(ctx, jobID, domain.JobStatusClosed))

		mock.Expect(mockExpectation{match: "update \"jobs\"", rowsAffected: 0})
		assert.ErrorIs(t, repo.UpdateStatus(ctx, jobID, domain.JobStatusClosed), gorm.ErrRecordNotFound)

		mock.Expect(mockExpectation{match: "update \"jobs\"", rowsAffected: 1})
		assert.NoError(t, repo.UpdateStatus(ctx, jobID, domain.JobStatusClosed))

		// Delete error, 0 rows, and success
		mock.Expect(mockExpectation{match: "delete from \"jobs\"", err: errors.New("db err")})
		assert.Error(t, repo.Delete(ctx, jobID))

		mock.Expect(mockExpectation{match: "delete from \"jobs\"", rowsAffected: 0})
		assert.ErrorIs(t, repo.Delete(ctx, jobID), gorm.ErrRecordNotFound)

		mock.Expect(mockExpectation{match: "delete from \"jobs\"", rowsAffected: 1})
		assert.NoError(t, repo.Delete(ctx, jobID))
	})
}
