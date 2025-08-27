package aws_s3

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
)

type AmazonS3Records struct {
	Records    filestore.Records
	bucketName string
	s3Bucket   *S3Bucket
}

func (s3Records *AmazonS3Records) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	err := s3Records.s3Bucket.UploadFile(ctx, s3Records.bucketName, model.Name.Slug, model.Url)
	if err != nil {
		return nil, err
	}
	return s3Records.Records.Add(ctx, model)
}
