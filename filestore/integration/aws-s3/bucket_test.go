package aws_s3

import (
	"context"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestBucket_UploadFile(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion("eu-north-1"))
	if err != nil {
		slog.Error(err.Error())
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := NewS3Bucket(s3Client, "naborlyjob")
	url, err := bucket.UploadFile(context.Background(), "bucket.go", "bucket.go")
	if err != nil {
		slog.Error(err.Error())
	}
	if len(url) <= 0 {
		slog.Error("Expected a valid url")
	}
}
