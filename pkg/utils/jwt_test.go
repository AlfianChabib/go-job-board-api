package utils_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/pkg/utils"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtManager(t *testing.T) {
	accessSecret := "test-access-secret-key"
	refreshSecret := "test-refresh-secret-key"
	manager := utils.NewJwtManager(accessSecret, 3600, refreshSecret, 7200)

	userID := uuid.New()
	role := domain.RoleCandidate

	t.Run("Generate and validate token pair successfully", func(t *testing.T) {
		pair, err := manager.GenerateTokenPair(userID, role)
		require.NoError(t, err)
		require.NotNil(t, pair)
		assert.NotEmpty(t, pair.AccessToken)
		assert.NotEmpty(t, pair.RefreshToken)
		assert.True(t, pair.AccessExpiredAt.After(time.Now()))
		assert.True(t, pair.RefreshExpiredAt.After(time.Now()))

		accessClaims, err := manager.ValidateAccessToken(pair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, accessClaims.UserID)
		assert.Equal(t, role, accessClaims.Role)

		refreshClaims, err := manager.ValidateRefreshToken(pair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, refreshClaims.UserID)
	})

	t.Run("Validate token with wrong secret fails", func(t *testing.T) {
		pair, err := manager.GenerateTokenPair(userID, role)
		require.NoError(t, err)

		// Validate access token using refresh token validator (different secret)
		_, err = manager.ValidateRefreshToken(pair.AccessToken)
		assert.Error(t, err)

		// Validate refresh token using access token validator (different secret)
		_, err = manager.ValidateAccessToken(pair.RefreshToken)
		assert.Error(t, err)
	})

	t.Run("Expired token returns Token has expired error", func(t *testing.T) {
		expiredManager := utils.NewJwtManager(accessSecret, -10, refreshSecret, -10)
		pair, err := expiredManager.GenerateTokenPair(userID, role)
		require.NoError(t, err)

		_, err = manager.ValidateAccessToken(pair.AccessToken)
		require.Error(t, err)
		assert.Equal(t, "Token has expired", err.Error())

		_, err = manager.ValidateRefreshToken(pair.RefreshToken)
		require.Error(t, err)
		assert.Equal(t, "Token has expired", err.Error())
	})

	t.Run("Unexpected signing method returns error", func(t *testing.T) {
		privKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		claims := &domain.JwtCustomClaims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		rsaToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privKey)
		require.NoError(t, err)

		_, err = manager.ValidateAccessToken(rsaToken)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unexpected signing token")
	})

	t.Run("Malformed token returns error", func(t *testing.T) {
		_, err := manager.ValidateAccessToken("not.a.valid.jwt")
		assert.Error(t, err)
	})
}
