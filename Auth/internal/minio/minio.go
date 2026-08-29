package miniocl

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(minioURL string) (*minio.Client, error) {
	minioClient, err := minio.New(minioURL, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadminpassword", ""), 
		Secure: false,                                                       
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect/authenticate with MinIO: %w", err)
	}

	return minioClient, nil
}
