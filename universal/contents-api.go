package universal

import "context"

// API

type Contents interface {
	ById(ctx context.Context, id int64) (Content, error)
	ByName(ctx context.Context, id string) (Content, error)
	Add(ctx context.Context, model *ContentModel) (Content, error)
}
