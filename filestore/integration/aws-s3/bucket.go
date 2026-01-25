package aws_s3

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Bucket struct {
	s3Client *s3.Client
	name     string
}

func NewS3Bucket(s3Client *s3.Client, name string) *S3Bucket {
	envName := os.Getenv("ENVIRONMENT")
	if envName == "" {
		envName = "dev"
	}
	return &S3Bucket{s3Client: s3Client, name: envName + "-" + name}
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics S3Bucket) UploadFile(ctx context.Context, objectKey string, pathToFile string) (string, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		slog.Error("Couldn't open file", "path", pathToFile, "error", err)
	} else {
		defer file.Close()
		_, err := basics.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(basics.name),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				slog.Error("Error while uploading object. The object is too large.\n"+
					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
					"or the multipart upload API (5TB max).", "bucket", basics.name)
			} else {
				slog.Error("Couldn't upload file", "path", pathToFile, "bucket", basics.name, "objectKey", objectKey, "error", err.Error())
			}
		} else {
			err = s3.NewObjectExistsWaiter(basics.s3Client).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(basics.name), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				slog.Error("Failed attempt to wait for object", "objectKey", objectKey, "error", err.Error())
			}
		}
		return basics.ObjectUrl(objectKey), err
	}
	return "", err
}

func (basics S3Bucket) PresignPutURL(
	ctx context.Context,
	objectKey string,
	expiration time.Duration,
	contentType string, // optional but recommended
) (string, error) {

	presigner := s3.NewPresignClient(basics.s3Client)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(basics.name),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}

	result, err := presigner.PresignPutObject(
		ctx,
		input,
		s3.WithPresignExpires(expiration),
	)
	if err != nil {
		slog.Error("Couldn't presign PUT URL",
			"bucket", basics.name,
			"objectKey", objectKey,
			"error", err,
		)
		return "", err
	}

	return result.URL, nil
}

func (basics S3Bucket) ObjectUrl(objectKey string) string {
	region := basics.s3Client.Options().Region
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", basics.name, region, objectKey)
}
