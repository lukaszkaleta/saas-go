package filestore

import (
	"context"
)

// API

type Records interface {
	Add(ctx context.Context, model *RecordModel) (Record, error)
	AddFromName(ctx context.Context, name string) (Record, error)
	AddFromUrl(ctx context.Context, url string) (Record, error)
	AddAll(ctx context.Context, model []*RecordModel) ([]Record, error)
	AddFromUrls(ctx context.Context, urls []string) ([]Record, error)
	Urls(ctx context.Context) ([]string, error)
}

type NoRecords struct {
}

func (r NoRecords) AddFromName(ctx context.Context, name string) (Record, error) {
	return nil, nil
}

func (r NoRecords) AddFromUrl(ctx context.Context, url string) (Record, error) {
	return nil, nil
}

func (r NoRecords) AddAll(ctx context.Context, model []*RecordModel) ([]Record, error) {
	return nil, nil
}

func (r NoRecords) AddFromUrls(ctx context.Context, urls []string) ([]Record, error) {
	return make([]Record, 0), nil
}

func (r NoRecords) Urls(ctx context.Context) ([]string, error) {
	return make([]string, 0), nil
}

func (NoRecords) Add(ctx context.Context, model *RecordModel) (Record, error) {
	return SolidRecord{model: model}, nil
}
