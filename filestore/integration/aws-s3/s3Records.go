package aws_s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type AmazonS3Records struct {
	records  filestore.Records
	s3Bucket *S3Bucket
}

func AmazonS3RecordsFromBucket(s3Bucket *S3Bucket, records filestore.Records) *AmazonS3Records {
	return &AmazonS3Records{records: records, s3Bucket: s3Bucket}
}

func AmazonS3RecordsFromClient(s3Client *s3.Client, bucketName string, records filestore.Records) *AmazonS3Records {
	return AmazonS3RecordsFromBucket(NewS3Bucket(s3Client, bucketName), records)
}

func (s3Records *AmazonS3Records) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	err := s3Records.s3Bucket.UploadFile(ctx, model.Name.Slug, model.Url)
	if err != nil {
		return nil, err
	}
	return s3Records.records.Add(ctx, model)
}
