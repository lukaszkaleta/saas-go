package category

type Categories interface {
	AddWithName(name string) (Category, error)
	AllLocalized(country string, language string) ([]*CategoryModel, error)
	ById(id int64) (Category, error)
}

func CategoryModels(Categories []Category) []*CategoryModel {
	var models []*CategoryModel
	for _, modelAware := range Categories {
		models = append(models, modelAware.Model()) // note the = instead of :=
	}
	return models
}
