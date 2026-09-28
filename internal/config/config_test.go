package config_test

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv(t *testing.T) {
	t.Run("LoadEnv reads .env file successfully", func(t *testing.T) {
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(origDir) }()

		tmpDir := t.TempDir()
		envContent := `APP_ENV=development
PORT=:8080
DATABASE_USER=postgres
DATABASE_PASSWORD=secret
DATABASE_DB=job_board
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_URL=postgres://postgres:secret@localhost:5432/job_board?sslmode=disable
ACCESS_SECRET_KEY=access-secret
ACCESS_DURATION=900
REFRESH_SECRET_KEY=refresh-secret
REFRESH_DURATION=604800
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_AVATAR_BUCKET=avatars
MINIO_CV_BUCKET=cvs
MINIO_USE_SSL=true
MINIO_ENDPOINT=localhost:9000
MINIO_PUBLIC_URL=http://localhost:9000
`
		err = os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(envContent), 0600)
		require.NoError(t, err)
		require.NoError(t, os.Chdir(tmpDir))

		env := config.LoadEnv()
		require.NotNil(t, env)
		assert.Equal(t, "development", env.AppEnv)
		assert.Equal(t, ":8080", env.Port)
		assert.Equal(t, "postgres", env.DatabaseUser)
		assert.Equal(t, "secret", env.DatabasePassword)
		assert.Equal(t, "job_board", env.DatabaseDB)
		assert.Equal(t, "localhost", env.DatabaseHost)
		assert.Equal(t, 5432, env.DatabasePort)
		assert.Equal(t, "postgres://postgres:secret@localhost:5432/job_board?sslmode=disable", env.DatabaseUrl)
		assert.Equal(t, "access-secret", env.AccessSecretKey)
		assert.Equal(t, 900, env.AccessDuration)
		assert.Equal(t, "refresh-secret", env.RefreshSecretKey)
		assert.Equal(t, 604800, env.RefreshDuration)
		assert.Equal(t, "minioadmin", env.MinioRootUser)
		assert.Equal(t, "minioadmin", env.MinioRootPassword)
		assert.Equal(t, "avatars", env.MinioAvatarBucket)
		assert.Equal(t, "cvs", env.MinioCvBucket)
		assert.True(t, env.MinioUseSSL)
		assert.Equal(t, "localhost:9000", env.MinioEndpoint)
		assert.Equal(t, "http://localhost:9000", env.MinioPublicUrl)
	})

	t.Run("LoadEnv panics when .env file is missing", func(t *testing.T) {
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(origDir) }()

		tmpDir := t.TempDir()
		require.NoError(t, os.Chdir(tmpDir))

		assert.Panics(t, func() {
			config.LoadEnv()
		})
	})
}
