package filestore

import (
	"context"
	"os"

	"github.com/lukaszkaleta/saas-go/universal"
)

func SingleDescriptionFromFile(ctx context.Context, records Records) universal.DescriptionFromFile {
	return func(file os.File) (*universal.DescriptionModel, error) {
		record, err := records.Add(ctx, &RecordModel{Url: file.Name()})
		if err != nil {
			return nil, err
		}
		return &universal.DescriptionModel{ImageUrl: record.Model().Url, Value: ""}, nil
	}
}
