package aws_s3

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

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
	url, err := bucket.UploadFile(t.Context(), "bucket.go", "bucket.go")
	if err != nil {
		slog.Error(err.Error())
	}
	if len(url) <= 0 {
		slog.Error("Expected a valid url")
	}
}

func TestBucket_PresignedUrl(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion("eu-north-1"))
	if err != nil {
		slog.Error(err.Error())
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := NewS3Bucket(s3Client, "naborly-prod-user-avatar")
	presignedURL, err := bucket.PresignPutURL(t.Context(), "presigned.png", time.Hour, "image/png")
	if err != nil {
		slog.Error(err.Error())
	}
	if len(presignedURL) <= 0 {
		slog.Error("Expected a valid url")
	}

	// Lets try to upload to presigned url.
	file, err := os.Open("presigned.png")
	if err != nil {
		slog.Error("Opening file failed")
		t.Fail()
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Read presigned png file failed")
		t.Fail()
	}

	req, err := http.NewRequest(
		http.MethodPut,
		presignedURL,
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		slog.Error("Create request failed")
		t.Fail()
	}

	// MUST match what was used when presigning
	req.Header.Set("Content-Type", "image/png")
	req.ContentLength = int64(len(fileBytes))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Send presigned png file failed")
		t.Fail()
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		slog.Info("upload failed:", "status", resp.StatusCode, "body", string(body))
		t.Fail()
	}
}
