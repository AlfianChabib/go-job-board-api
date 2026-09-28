package repository_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("Register success and error", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewAuthRepository(db)

		// 1. Error
		mock.Expect(mockExpectation{
			match: "insert into \"users\"",
			err:   errors.New("insert failed"),
		})
		res, err := repo.Register(ctx, domain.User{Name: "John", Email: "john@example.com"})
		assert.Nil(t, res)
		assert.EqualError(t, err, "Failed to register new user")

		// 2. Success
		mock.Expect(mockExpectation{
			match:        "insert into \"users\"",
			rowsAffected: 1,
		})
		res, err = repo.Register(ctx, domain.User{Name: "John", Email: "john@example.com"})
		require.NoError(t, err)
		assert.Equal(t, "john@example.com", res.Email)
		assert.NotEqual(t, uuid.Nil, res.ID)
	})

	t.Run("IsUserExist true and false", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewAuthRepository(db)

		// 1. Not found / error -> false
		mock.Expect(mockExpectation{
			match: "select \"email\" from \"users\"",
			err:   errors.New("record not found"),
		})
		assert.False(t, repo.IsUserExist(ctx, "missing@example.com"))

		// 2. Found -> true
		mock.Expect(mockExpectation{
			match:   "select \"email\" from \"users\"",
			columns: []string{"email"},
			rows:    [][]driver.Value{{"found@example.com"}},
		})
		assert.True(t, repo.IsUserExist(ctx, "found@example.com"))
	})

	t.Run("FindByEmail success and error", func(t *testing.T) {
		db, mock := newMockGormDB(t)
		repo := repository.NewAuthRepository(db)

		// 1. Error
		mock.Expect(mockExpectation{
			match: "select * from \"users\"",
			err:   errors.New("not found"),
		})
		res, err := repo.FindByEmail(ctx, "missing@example.com")
		assert.Nil(t, res)
		assert.EqualError(t, err, "User not found")

		// 2. Success with Preload("Auth")
		userID := uuid.New()
		authID := uuid.New()
		mock.Expect(mockExpectation{
			match:   "select * from \"users\"",
			columns: []string{"id", "name", "email", "role"},
			rows:    [][]driver.Value{{userID.String(), "John", "john@example.com", "CANDIDATE"}},
		})
		mock.Expect(mockExpectation{
			match:   "select * from \"auths\"",
			columns: []string{"id", "user_id", "email", "password"},
			rows:    [][]driver.Value{{authID.String(), userID.String(), "john@example.com", "hashed-pw"}},
		})

		res, err = repo.FindByEmail(ctx, "john@example.com")
		require.NoError(t, err)
		assert.Equal(t, userID, res.ID)
		assert.Equal(t, "hashed-pw", res.Auth.Password)
	})
}
