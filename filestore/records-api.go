package filestore

import (
	"context"
)

// API

type Records interface {
	Add(ctx context.Context, model *RecordModel) (Record, error)
}
