package repository

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

type storageRepository struct {
	store          *minio.Client
	minioPublicUrl string
	avatarBucket   string
	cvBucket       string
}

func NewStorageRepository(store *minio.Client, minioPublicUrl string, avatarBucket string, cvBucket string) StorageRepository {
	return &storageRepository{
		store:          store,
		minioPublicUrl: minioPublicUrl,
		avatarBucket:   avatarBucket,
		cvBucket:       cvBucket,
	}
}

func (repo *storageRepository) UploadAvatar(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (*string, error) {
	putOpts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	uploadInfo, err := repo.store.PutObject(
		ctx,
		repo.avatarBucket,
		fileName,
		file,
		size,
		putOpts,
	)
	if err != nil {
		return nil, err
	}

	log.Printf("file uploaded successfully: ETag: %s, Size: %d", uploadInfo.ETag, uploadInfo.Size)

	url := fmt.Sprintf("%s/%s/%s", repo.minioPublicUrl, repo.avatarBucket, fileName)

	return &url, nil
}

func (repo *storageRepository) DeleteAvatar(ctx context.Context, fileName string) error {
	panic("TODO: Implement")
}
