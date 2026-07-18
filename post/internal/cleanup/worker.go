package cleanup

import (
	"context"


	"github.com/abelmalu/golang-posts/platform"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type CleanupJob struct {
	Bucket string
	Object string
}

type CleanUpService struct {
	minioClient *minio.Client
	logger      *platform.Logger
	jobs        chan CleanupJob
}

func NewCleanUpService(minioClient *minio.Client, logger *platform.Logger) *CleanUpService {

	cleanupService := &CleanUpService{

		minioClient: minioClient,
		logger:      logger,
		jobs:        make(chan CleanupJob),
	}

	for range 5 {

		go cleanupService.workers()

	}

	return cleanupService

}

func (cleanupService CleanUpService) workers() {

	for job := range cleanupService.jobs {

		err := cleanupService.minioClient.RemoveObject(
			context.Background(),
			job.Bucket,
			job.Object,
			minio.RemoveObjectOptions{},
		)

		if err != nil {
			cleanupService.logger.Error(
				"failed to delete object",
				zap.Error(err),
			)
		}

	}
}




func (cleanupService *CleanUpService) DeleteObject(bucket, object string) {

    select {
    case cleanupService.jobs <- CleanupJob{
        Bucket: bucket,
        Object: object,
    }:

    default:
        cleanupService.logger.Warn("cleanup queue full")
    }
}
