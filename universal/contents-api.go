package universal

import "context"

// API

type Contents interface {
	ById(ctx context.Context, id int64) (Content, error)
	ByName(ctx context.Context, id string) (Content, error)
	Add(ctx context.Context, model *ContentModel) (Content, error)
}

// Builders

func NewSolidContents(ratings Contents) SolidContents {
	return SolidContents{ratings: ratings}
}

// Solid

type SolidContents struct {
	ratings Contents
}

func (s SolidContents) Add(ctx context.Context, r ContentModel) (*Content, error) {
	return nil, nil
}

func (s SolidContents) ById(ctx context.Context, id int64) (Content, error) {
	return s.ratings.ById(ctx, id)
}
func (s SolidContents) ByName(ctx context.Context, name string) (Content, error) {
	return s.ratings.ByName(ctx, name)
}
