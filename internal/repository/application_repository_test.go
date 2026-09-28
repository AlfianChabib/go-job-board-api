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
)

func TestApplicationRepository(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()
	jobID := uuid.New()
	candidateID := uuid.New()
	compID := uuid.New()
	profileID := uuid.New()

	t.Run("Create, FindById, and FindByJobAndCandidate", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewApplicationRepository(db)

		// Create error & success
		mock.Expect(mockExpectation{match: "insert into \"applications\"", err: errors.New("insert err")})
		_, err := repo.Create(ctx, domain.Application{JobId: jobID, CandidateId: candidateID})
		assert.Error(t, err)

		mock.Expect(mockExpectation{match: "insert into \"applications\"", rowsAffected: 1})
		created, err := repo.Create(ctx, domain.Application{JobId: jobID, CandidateId: candidateID})
		require.NoError(t, err)
		assert.Equal(t, jobID, created.JobId)

		// FindById error & success
		mock.Expect(mockExpectation{match: "select * from \"applications\"", err: errors.New("not found")})
		_, err = repo.FindById(ctx, appID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id", "status"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String(), "APPLIED"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"users\"",
			columns: []string{"id", "name", "email"},
			rows:    [][]driver.Value{{candidateID.String(), "Alice", "alice@example.com"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), candidateID.String()}},
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
		found, err := repo.FindById(ctx, appID)
		require.NoError(t, err)
		assert.Equal(t, appID, found.ID)

		// FindByJobAndCandidate error & success
		mock.Expect(mockExpectation{match: "select * from \"applications\"", err: errors.New("not found")})
		_, err = repo.FindByJobAndCandidate(ctx, jobID, candidateID)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String()}},
		})
		found2, err := repo.FindByJobAndCandidate(ctx, jobID, candidateID)
		require.NoError(t, err)
		assert.Equal(t, appID, found2.ID)
	})

	t.Run("FindAllByJobId and FindAllByCandidateId", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewApplicationRepository(db)

		// FindAllByJobId count error, find error, and success
		jobReq := web.GetJobApplicationsRequest{Status: "APPLIED", Search: "alice"}
		mock.Expect(mockExpectation{match: "select count(*)", err: errors.New("count err")})
		_, _, err := repo.FindAllByJobId(ctx, jobID, jobReq)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{match: "select applications.*", err: errors.New("find err")})
		_, _, err = repo.FindAllByJobId(ctx, jobID, jobReq)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match:   "select applications.*",
			columns: []string{"id", "job_id", "candidate_id"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"users\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{candidateID.String(), "Alice"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"profiles\"",
			columns: []string{"id", "user_id"},
			rows:    [][]driver.Value{{profileID.String(), candidateID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"jobs\"",
			columns: []string{"id", "company_id"},
			rows:    [][]driver.Value{{jobID.String(), compID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		apps, total, err := repo.FindAllByJobId(ctx, jobID, jobReq)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, apps, 1)

		// FindAllByCandidateId count error, find error, and success
		candReq := web.GetCandidateApplicationsRequest{Status: "APPLIED"}
		mock.Expect(mockExpectation{match: "select count(*)", err: errors.New("count err")})
		_, _, err = repo.FindAllByCandidateId(ctx, candidateID, candReq)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{match: "select * from \"applications\"", err: errors.New("find err")})
		_, _, err = repo.FindAllByCandidateId(ctx, candidateID, candReq)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select count(*)",
			columns: []string{"count"},
			rows:    [][]driver.Value{{int64(1)}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"jobs\"",
			columns: []string{"id", "company_id"},
			rows:    [][]driver.Value{{jobID.String(), compID.String()}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"companies\"",
			columns: []string{"id", "name"},
			rows:    [][]driver.Value{{compID.String(), "Acme"}},
		})
		apps, total, err = repo.FindAllByCandidateId(ctx, candidateID, candReq)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, apps, 1)
	})

	t.Run("UpdateStatus and Withdraw", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewApplicationRepository(db)
		notes := "Great fit"
		reason := "No longer available"

		// UpdateStatus not found, save error, and success
		mock.Expect(mockExpectation{match: "select * from \"applications\"", err: errors.New("not found")})
		_, err := repo.UpdateStatus(ctx, appID, domain.ApplicationStatusShortlisted, &notes, &reason)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id", "status"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String(), "APPLIED"}},
		})
		mock.Expect(mockExpectation{match: "select * from \"users\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "update \"applications\"", err: errors.New("save err")})
		_, err = repo.UpdateStatus(ctx, appID, domain.ApplicationStatusShortlisted, &notes, &reason)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id", "status"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String(), "APPLIED"}},
		})
		mock.Expect(mockExpectation{match: "select * from \"users\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "update \"applications\"", rowsAffected: 1})
		updated, err := repo.UpdateStatus(ctx, appID, domain.ApplicationStatusShortlisted, &notes, &reason)
		require.NoError(t, err)
		assert.Equal(t, domain.ApplicationStatusShortlisted, updated.Status)

		// Withdraw not found, save error, and success
		mock.Expect(mockExpectation{match: "select * from \"applications\"", err: errors.New("not found")})
		_, err = repo.Withdraw(ctx, appID, &reason)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id", "status"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String(), "APPLIED"}},
		})
		mock.Expect(mockExpectation{match: "select * from \"users\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "update \"applications\"", err: errors.New("save err")})
		_, err = repo.Withdraw(ctx, appID, &reason)
		assert.Error(t, err)

		mock.Expect(mockExpectation{
			match:   "select * from \"applications\"",
			columns: []string{"id", "job_id", "candidate_id", "status"},
			rows:    [][]driver.Value{{appID.String(), jobID.String(), candidateID.String(), "APPLIED"}},
		})
		mock.Expect(mockExpectation{match: "select * from \"users\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "select * from \"jobs\"", columns: []string{"id"}, rows: nil})
		mock.Expect(mockExpectation{match: "update \"applications\"", rowsAffected: 1})
		withdrawn, err := repo.Withdraw(ctx, appID, &reason)
		require.NoError(t, err)
		assert.Equal(t, domain.ApplicationStatusWithdrawn, withdrawn.Status)
	})
}
