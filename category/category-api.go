package category

import "github.com/lukaszkaleta/saas-go/universal"

// API

type Category interface {
	Model() *CategoryModel
	Update(newModel *CategoryModel) error
	Localizations() universal.Localizations
}

// Builder

func EmptyCategoryModel() *CategoryModel {
	return &CategoryModel{
		Id:          0,
		ParentId:    nil,
		Name:        universal.EmptyNameModel(),
		Description: universal.EmptyDescriptionModel(),
	}
}

// Model

type CategoryModel struct {
	Id          int64                       `json:"id"`
	ParentId    *int64                      `json:"parentId"`
	Description *universal.DescriptionModel `json:"description"`
	Name        *universal.NameModel        `json:"name"`
}

func (model *CategoryModel) Change(newModel *CategoryModel) {
	model.Description.Change(newModel.Description)
	model.Name.Change(newModel.Name.Value)
}

func (model *CategoryModel) GetId() int64 {
	return model.Id
}

// Solid

type SolidCategory struct {
	model    *CategoryModel
	Category Category
}

func NewSolidCategory(model *CategoryModel, category Category) Category {
	return &SolidCategory{model, category}
}

func (category SolidCategory) Update(newModel *CategoryModel) error {
	category.model.Change(newModel)
	if category.Category == nil {
		return nil
	}
	return category.Category.Update(newModel)
}

func (category SolidCategory) Model() *CategoryModel {
	return category.model
}

func (category SolidCategory) Localizations() universal.Localizations {
	return category.Category.Localizations()
}
