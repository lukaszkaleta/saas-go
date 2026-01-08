package universal

import "context"

// API

type Content interface {
	Model(ctx context.Context) *ContentModel
	Update(ctx context.Context, newModel *ContentModel) error
	Localizations() Localizations
}

// Builder

func EmptyContentModel() *ContentModel {
	return &ContentModel{
		Id:    0,
		Name:  EmptyNameModel(),
		Value: "{}",
	}
}

// Model

type ContentModel struct {
	Id    int64      `json:"id"`
	Name  *NameModel `json:"name"`
	Value string     `json:"value"`
}

func (model *ContentModel) Change(newModel *ContentModel) {
	model.Name.Change(newModel.Name.Value)
	model.Value = newModel.Value
}

// Solid

type SolidContent struct {
	model   *ContentModel
	Content Content
}

func NewSolidContent(model *ContentModel, category Content) Content {
	return &SolidContent{model, category}
}

func (category SolidContent) Update(ctx context.Context, newModel *ContentModel) error {
	category.model.Change(newModel)
	if category.Content == nil {
		return nil
	}
	return category.Content.Update(ctx, newModel)
}

func (category SolidContent) Model(ctx context.Context) *ContentModel {
	return category.model
}

func (category SolidContent) Localizations() Localizations {
	return category.Content.Localizations()
}
