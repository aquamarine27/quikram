package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client *minio.Client
	bucket string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucket, region string) *S3Storage {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
		Region: region,
	})
	if err != nil {
		log.Fatalf("failed to create S3 client: %v", err)
	}

	ctx := context.Background()
	if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region}); err != nil {
		exists, errExists := client.BucketExists(ctx, bucket)
		if errExists != nil || !exists {
			log.Fatalf("failed to create bucket %s: %v", bucket, err)
		}
		log.Printf("using existing S3 bucket: %s", bucket)
	} else {
		log.Printf("created S3 bucket: %s", bucket)
	}

	return &S3Storage{client: client, bucket: bucket}
}

func (s *S3Storage) Upload(key string, reader io.Reader, contentType string) error {
	_, err := s.client.PutObject(context.Background(), s.bucket, key, reader, -1,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}
	return nil
}

func (s *S3Storage) Download(key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(context.Background(), s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("s3 download failed: %w", err)
	}
	return obj, nil
}

func (s *S3Storage) Delete(key string) error {
	return s.client.RemoveObject(context.Background(), s.bucket, key, minio.RemoveObjectOptions{})
}
