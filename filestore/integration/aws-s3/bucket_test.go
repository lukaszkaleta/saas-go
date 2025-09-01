package aws_s3

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestBucket_UploadFile(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := NewS3Bucket(s3Client, "naborly-dev-offer")
	url, err := bucket.UploadFile(context.Background(), "bucket.go", "bucket.go")
	if err != nil {
		log.Fatal(err)
	}
	if len(url) <= 0 {
		log.Fatal("Expected a valid url")
	}
}
