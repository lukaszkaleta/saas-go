package aws_s3

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

func GenerateS3Description(ctx context.Context, s3Bucket *S3Bucket, objectKey string, description universal.Description) (string, error) {
	presignedUrl, err := s3Bucket.PresignPutURL(ctx, objectKey, time.Hour, "image/jpg")
	if err != nil {
		return "", err
	}
	url := s3Bucket.ObjectUrl(objectKey)
	err = description.UpdateImageUrl(ctx, &url)
	if err != nil {
		return presignedUrl, err
	}
	return presignedUrl, err
}
