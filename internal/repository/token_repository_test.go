package repository_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTokenRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("Save success and error", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewTokenRepository(db)

		// 1. Error
		mock.Expect(mockExpectation{
			match: "tokens",
			err:   errors.New("db error"),
		})
		res, err := repo.Save(ctx, domain.Token{RefreshToken: "tok"})
		assert.Nil(t, res)
		assert.Error(t, err)

		// 2. Success
		mock.Expect(mockExpectation{
			match:        "tokens",
			rowsAffected: 1,
		})
		res, err = repo.Save(ctx, domain.Token{RefreshToken: "tok", ExpiresAt: time.Now().Add(time.Hour)})
		require.NoError(t, err)
		assert.Equal(t, "tok", res.RefreshToken)
	})

	t.Run("RevokeToken db error, 0 rows affected, and success", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewTokenRepository(db)

		// 1. DB error
		mock.Expect(mockExpectation{
			match: "update \"tokens\"",
			err:   errors.New("db failure"),
		})
		err := repo.RevokeToken(ctx, "tok")
		assert.EqualError(t, err, "Internal server error")

		// 2. 0 rows affected
		mock.Expect(mockExpectation{
			match:        "update \"tokens\"",
			rowsAffected: 0,
		})
		err = repo.RevokeToken(ctx, "tok")
		assert.EqualError(t, err, "Sesi tidak valid atau sudah berakhir")

		// 3. Success
		mock.Expect(mockExpectation{
			match:        "update \"tokens\"",
			rowsAffected: 1,
		})
		err = repo.RevokeToken(ctx, "tok")
		assert.NoError(t, err)
	})

	t.Run("FindTokenWithUser not found, generic error, and success", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewTokenRepository(db)
		userID := uuid.New()
		tokID := uuid.New()

		// 1. Record not found
		mock.Expect(mockExpectation{
			match: "select * from \"tokens\"",
			err:   gorm.ErrRecordNotFound,
		})
		res, err := repo.FindTokenWithUser(ctx, userID, "tok")
		assert.Nil(t, res)
		assert.EqualError(t, err, "token tidak valid atau sudah dicabut")

		// 2. Generic error
		mock.Expect(mockExpectation{
			match: "select * from \"tokens\"",
			err:   errors.New("connection lost"),
		})
		res, err = repo.FindTokenWithUser(ctx, userID, "tok")
		assert.Nil(t, res)
		assert.EqualError(t, err, "connection lost")

		// 3. Success with Preload("User")
		mock.Expect(mockExpectation{
			match:   "select * from \"tokens\"",
			columns: []string{"id", "user_id", "refresh_token"},
			rows:    [][]driver.Value{{tokID.String(), userID.String(), "tok"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"users\"",
			columns: []string{"id", "name", "email", "role"},
			rows:    [][]driver.Value{{userID.String(), "Alice", "alice@example.com", "CANDIDATE"}},
		})
		res, err = repo.FindTokenWithUser(ctx, userID, "tok")
		require.NoError(t, err)
		assert.Equal(t, tokID, res.ID)
		assert.Equal(t, userID, res.User.ID)
	})
}
