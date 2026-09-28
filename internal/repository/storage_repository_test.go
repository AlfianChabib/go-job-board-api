package repository_test

import (
	"AlfianChabib/go-job-board-api/internal/repository"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestMinioClient(t *testing.T, shouldFail bool) (*minio.Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, hasLoc := r.URL.Query()["location"]; hasLoc {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`))
			return
		}

		if shouldFail {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>AccessDenied</Code><Message>Access Denied</Message></Error>`))
			return
		}

		switch r.Method {
		case http.MethodPut:
			w.Header().Set("ETag", `"etag-12345"`)
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))

	endpoint := strings.TrimPrefix(server.URL, "http://")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)
	return client, server
}

func TestStorageRepository(t *testing.T) {
	ctx := context.Background()
	publicURL := "http://localhost:9000"
	avatarBucket := "avatars"
	cvBucket := "cvs"

	t.Run("UploadAvatar, UploadLogo, UploadBanner, and DeleteAvatar success", func(t *testing.T) {
		client, server := createTestMinioClient(t, false)
		defer server.Close()

		repo := repository.NewStorageRepository(client, publicURL, avatarBucket, cvBucket)
		content := []byte("fake-image-bytes")

		avatarURL, err := repo.UploadAvatar(ctx, "avatar_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		require.NoError(t, err)
		require.NotNil(t, avatarURL)
		assert.Equal(t, "http://localhost:9000/avatars/avatar_1.png", *avatarURL)

		logoURL, err := repo.UploadLogo(ctx, "logo_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		require.NoError(t, err)
		require.NotNil(t, logoURL)
		assert.Equal(t, "http://localhost:9000/avatars/logo_1.png", *logoURL)

		bannerURL, err := repo.UploadBanner(ctx, "banner_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		require.NoError(t, err)
		require.NotNil(t, bannerURL)
		assert.Equal(t, "http://localhost:9000/avatars/banner_1.png", *bannerURL)

		err = repo.DeleteAvatar(ctx, "avatar_1.png")
		require.NoError(t, err)
	})

	t.Run("UploadAvatar, UploadLogo, UploadBanner, and DeleteAvatar errors", func(t *testing.T) {
		client, server := createTestMinioClient(t, true)
		defer server.Close()

		repo := repository.NewStorageRepository(client, publicURL, avatarBucket, cvBucket)
		content := []byte("fake-image-bytes")

		_, err := repo.UploadAvatar(ctx, "avatar_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		assert.Error(t, err)

		_, err = repo.UploadLogo(ctx, "logo_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		assert.Error(t, err)

		_, err = repo.UploadBanner(ctx, "banner_1.png", bytes.NewReader(content), int64(len(content)), "image/png")
		assert.Error(t, err)

		err = repo.DeleteAvatar(ctx, "avatar_1.png")
		assert.Error(t, err)
	})
}
