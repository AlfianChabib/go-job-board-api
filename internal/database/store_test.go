package database

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockMinioServer(existingBuckets map[string]bool) (*httptest.Server, *sync.Map) {
	createdBuckets := &sync.Map{}
	for k, v := range existingBuckets {
		if v {
			createdBuckets.Store(k, true)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.SplitN(path, "/", 2)
		bucket := parts[0]

		if _, hasLoc := r.URL.Query()["location"]; hasLoc {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`))
			return
		}

		if _, hasPolicy := r.URL.Query()["policy"]; hasPolicy && r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		switch r.Method {
		case http.MethodHead:
			if _, ok := createdBuckets.Load(bucket); ok {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPut:
			createdBuckets.Store(bucket, true)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))

	return server, createdBuckets
}

func TestOpenMinioClientAndSetupBucket(t *testing.T) {
	t.Run("OpenMinioClient creates non-existing public and private buckets", func(t *testing.T) {
		server, createdBuckets := newMockMinioServer(nil)
		defer server.Close()

		endpoint := strings.TrimPrefix(server.URL, "http://")
		env := &config.Env{
			MinioEndpoint:     endpoint,
			MinioRootUser:     "minioadmin",
			MinioRootPassword: "minioadmin",
			MinioUseSSL:       false,
			MinioAvatarBucket: "avatars",
			MinioCvBucket:     "cvs",
		}

		client := OpenMinioClient(env)
		require.NotNil(t, client)

		_, okAvatar := createdBuckets.Load("avatars")
		_, okCv := createdBuckets.Load("cvs")
		assert.True(t, okAvatar)
		assert.True(t, okCv)
	})

	t.Run("setupBucket skips creation when bucket already exists", func(t *testing.T) {
		server, _ := newMockMinioServer(map[string]bool{"existing-bucket": true})
		defer server.Close()

		endpoint := strings.TrimPrefix(server.URL, "http://")
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		require.NoError(t, err)

		setupBucket(context.Background(), client, "existing-bucket", true)
	})
}
