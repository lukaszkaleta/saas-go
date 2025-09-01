package universal

import (
	"os"
)

// API

type Description interface {
	Model() *DescriptionModel
	Update(newModel *DescriptionModel) error
	UpdateImageUrl(imageUrl *string) error
}

// Builder

func DescriptionFromModel(model *DescriptionModel) Description {
	return SolidDescription{
		model:       model,
		Description: nil,
	}
}

// Model

type DescriptionModel struct {
	Value    string `json:"value"`
	ImageUrl string `json:"imageUrl"`
}

func EmptyDescriptionModel() *DescriptionModel {
	return &DescriptionModel{
		Value:    "",
		ImageUrl: "",
	}
}

func (model *DescriptionModel) Change(newModel *DescriptionModel) {
	model.Value = newModel.Value
	model.ImageUrl = newModel.ImageUrl
}

// Solid

type SolidDescription struct {
	model       *DescriptionModel
	Description Description
}

func NewSolidDescription(model *DescriptionModel, Description Description) Description {
	return &SolidDescription{model, Description}
}

func (sd SolidDescription) Update(newModel *DescriptionModel) error {
	sd.model.Change(newModel)
	if sd.Description == nil {
		return nil
	}
	return sd.Description.Update(newModel)
}

func (sd SolidDescription) UpdateImageUrl(imageUrl *string) error {
	sd.model.ImageUrl = *imageUrl
	if sd.Description == nil {
		return nil
	}
	return sd.Description.UpdateImageUrl(imageUrl)
}

func (sd SolidDescription) Model() *DescriptionModel {
	return sd.model
}

type DescriptionFromUrl func(url string) (*DescriptionModel, error)
type DescriptionFromFile func(file os.File) (*DescriptionModel, error)
