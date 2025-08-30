package universal

import (
	"os"
)

// API

type Description interface {
	Model() *DescriptionModel
	Update(newModel *DescriptionModel) error
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

func (addr SolidDescription) Update(newModel *DescriptionModel) error {
	addr.model.Change(newModel)
	if addr.Description == nil {
		return nil
	}
	return addr.Description.Update(newModel)
}

func (addr SolidDescription) Model() *DescriptionModel {
	return addr.model
}

type DescriptionFromUrl func(url string) (*DescriptionModel, error)
type DescriptionFromFile func(file os.File) (*DescriptionModel, error)
