package repository

import (
	"context"
	"io"
)

type StorageRepository interface {
	UploadAvatar(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (*string, error)
	DeleteAvatar(ctx context.Context, fileName string) error
}
