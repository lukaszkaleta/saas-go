package universal

import "context"

// API

type Localization interface {
	Model(ctx context.Context) *LocalizationModel
	Update(ctx context.Context, newModel *LocalizationModel) error
}

// Builder

func LocalizationFromModel(model *LocalizationModel) Localization {
	return SolidLocalization{
		model:        model,
		Localization: nil,
	}
}

// Model

type LocalizationModel struct {
	Id          int64      `json:"id"`
	OwnerId     int64      `json:"owner_id"`
	Country     string     `json:"country"`
	Language    string     `json:"language"`
	Translation *NameModel `json:"translation"`
}

func (model *LocalizationModel) Change(newModel *LocalizationModel) {
	model.Country = newModel.Country
	model.Language = newModel.Language
	model.Translation.Change(newModel.Translation.Value)
}

// Solid

type SolidLocalization struct {
	model        *LocalizationModel
	Localization Localization
}

func NewSolidLocalization(model *LocalizationModel, Localization Localization) Localization {
	return &SolidLocalization{model, Localization}
}

func (addr SolidLocalization) Update(ctx context.Context, newModel *LocalizationModel) error {
	addr.model.Change(newModel)
	if addr.Localization == nil {
		return nil
	}
	return addr.Localization.Update(ctx, newModel)
}

func (addr SolidLocalization) Model(ctx context.Context) *LocalizationModel {
	return addr.model
}
