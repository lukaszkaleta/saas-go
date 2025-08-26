package universal

// API

type Tag interface {
	Model() *TagModel
	Update(newModel *TagModel) error
}

// Builder

func TagFromModel(model *TagModel) Tag {
	return SolidTag{
		model: model,
	}
}

// Model

type TagModel struct {
}

func (model *TagModel) Change(newModel *TagModel) {

}

// Solid

type SolidTag struct {
	model *TagModel
	Tag   Tag
}

func NewSolidTag(model *TagModel, Tag Tag) Tag {
	return &SolidTag{model, Tag}
}

func (Tag SolidTag) Update(newModel *TagModel) error {
	Tag.model.Change(newModel)
	if Tag.Tag == nil {
		return nil
	}
	return Tag.Tag.Update(newModel)
}

func (Tag SolidTag) Model() *TagModel {
	return Tag.model
}
