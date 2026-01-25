package aws_s3

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

func GenerateS3Description(ctx context.Context, s3Bucket *S3Bucket, description universal.Description) (string, error) {
	model := description.Model(ctx)
	objectKey := model.Value
	presignedUrl, err := s3Bucket.PresignPutURL(ctx, objectKey, time.Hour, "image/png")
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
