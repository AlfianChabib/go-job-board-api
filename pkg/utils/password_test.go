package utils_test

import (
	"AlfianChabib/go-job-board-api/pkg/utils"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher(t *testing.T) {
	hasher := utils.NewBcryptHasher(bcrypt.MinCost)

	t.Run("Hash and Compare valid password", func(t *testing.T) {
		password := []byte("SuperSecret123!")
		hashed, err := hasher.Hash(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, string(password), hashed)

		assert.True(t, hasher.Compare([]byte(hashed), password))
		assert.False(t, hasher.Compare([]byte(hashed), []byte("WrongPassword")))
	})

	t.Run("Compare with empty hash or password returns false", func(t *testing.T) {
		assert.False(t, hasher.Compare([]byte(""), []byte("password")))
		assert.False(t, hasher.Compare(nil, []byte("password")))
		assert.False(t, hasher.Compare([]byte("somehash"), []byte("")))
		assert.False(t, hasher.Compare([]byte("somehash"), nil))
	})

	t.Run("Hash error when password exceeds 72 bytes or invalid cost", func(t *testing.T) {
		longPassword := bytes.Repeat([]byte("a"), 80)
		_, err := hasher.Hash(longPassword)
		assert.Error(t, err)

		invalidCostHasher := utils.NewBcryptHasher(100)
		_, err = invalidCostHasher.Hash([]byte("short"))
		assert.Error(t, err)
	})
}
