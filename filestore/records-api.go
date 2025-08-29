package filestore

import (
	"context"
)

// API

type Records interface {
	Add(ctx context.Context, model *RecordModel) (Record, error)
}

type NoRecords struct {
}

func (NoRecords) Add(ctx context.Context, model *RecordModel) (Record, error) {
	return SolidRecord{model: model}, nil
}
