package category

import "context"

type Categories interface {
	AddWithName(ctx context.Context, name string) (Category, error)
	AllLocalized(ctx context.Context, country string, language string) ([]*CategoryModel, error)
	ById(ctx context.Context, id int64) (Category, error)
	ByIds(ctx context.Context, id []int64) ([]CategoryModel, error)
}

func Models(ctx context.Context, Categories []Category) []*CategoryModel {
	var models []*CategoryModel
	for _, modelAware := range Categories {
		models = append(models, modelAware.Model(ctx)) // note the = instead of :=
	}
	return models
}
