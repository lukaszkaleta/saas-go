package universal

import "context"

// API

type Position interface {
	Model(ctx context.Context) *PositionModel
	Update(ctx context.Context, newModel *PositionModel) error
}

// Builder

func PositionFromModel(model *PositionModel) Position {
	return SolidPosition{
		model: model,
	}
}

func NewPosition(lat float64, lon float64) Position {
	return PositionFromModel(&PositionModel{lat, lon})
}

// Model

type PositionModel struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func EmptyPositionModel() *PositionModel {
	return &PositionModel{}
}

func (model *PositionModel) Change(newModel *PositionModel) {
	model.Lat = newModel.Lat
	model.Lon = newModel.Lon
}

// Solid

type SolidPosition struct {
	model    *PositionModel
	position Position
}

func NewSolidPosition(model *PositionModel, Position Position) Position {
	return &SolidPosition{model, Position}
}

func (addr SolidPosition) Update(ctx context.Context, newModel *PositionModel) error {
	addr.model.Change(newModel)
	if addr.position == nil {
		return nil
	}
	return addr.position.Update(ctx, newModel)
}

func (addr SolidPosition) Model(ctx context.Context) *PositionModel {
	return addr.model
}
