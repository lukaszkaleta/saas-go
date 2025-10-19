package universal

// API

type Content interface {
	Model() *ContentModel
	Update(newModel *ContentModel) error
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

func (category SolidContent) Update(newModel *ContentModel) error {
	category.model.Change(newModel)
	if category.Content == nil {
		return nil
	}
	return category.Content.Update(newModel)
}

func (category SolidContent) Model() *ContentModel {
	return category.model
}

func (category SolidContent) Localizations() Localizations {
	return category.Content.Localizations()
}
