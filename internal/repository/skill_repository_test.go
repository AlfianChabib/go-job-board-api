package repository_test

import (
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

func TestSkillRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("FindAll with default limit, search, and error", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewSkillRepository(db)

		// 1. Error
		mock.Expect(mockExpectation{
			match: "select * from \"skills\"",
			err:   errors.New("db error"),
		})
		res, err := repo.FindAll(ctx, web.GetDataRequest{Search: "go", Limit: 0})
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")

		// 2. Success with search and custom limit
		skillID := uuid.New()
		mock.Expect(mockExpectation{
			match:   "select * from \"skills\"",
			columns: []string{"id", "name", "label", "abbreviation"},
			rows:    [][]driver.Value{{skillID.String(), "golang", "Go", "Go"}},
		})
		res, err = repo.FindAll(ctx, web.GetDataRequest{Search: "go", Limit: 5})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, skillID, res[0].ID)
		assert.Equal(t, "golang", res[0].Name)
	})
}
