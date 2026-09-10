package database

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// OpenStore mengembalikan *minio.Client resmi, bukan gofiber/storage
func OpenMinioClient(env *config.Env) *minio.Client {
	ctx := context.Background()

	// 1. Inisialisasi Koneksi ke Server MinIO (Client level, bukan Bucket level)
	client, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioRootUser, env.MinioRootPassword, ""),
		Secure: env.MinioUseSSL,
	})
	if err != nil {
		log.Fatalf("failed to connect to minio: %v", err)
	}

	// 2. Inisialisasi Bucket Avatar (Diatur menjadi Public)
	setupBucket(ctx, client, env.MinioAvatarBucket, true)

	// 3. Inisialisasi Bucket CV (Diatur menjadi Private)
	setupBucket(ctx, client, env.MinioCvBucket, false)

	return client
}

// setupBucket adalah helper private untuk mengecek, membuat, dan mengatur policy
func setupBucket(ctx context.Context, client *minio.Client, bucketName string, isPublic bool) {
	// Cek ketersediaan bucket
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("failed to check bucket %s: %v", bucketName, err)
	}

	// Jika belum ada, buat bucket
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket %s: %v", bucketName, err)
		}
		log.Printf("Bucket '%s' created successfully.", bucketName)

		// Terapkan Public Read-Only Policy jika bucket digunakan untuk Avatar
		if isPublic {
			policy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::%s/*"]
					}
				]
			}`, bucketName)

			err = client.SetBucketPolicy(ctx, bucketName, policy)
			if err != nil {
				log.Fatalf("failed to set public policy for bucket %s: %v", bucketName, err)
			}
			log.Printf("Bucket '%s' is now public.", bucketName)
		}
	}
}
