package universal

// API

type Position interface {
	Model() *PositionModel
	Update(newModel *PositionModel) error
}

// Builder

func PositionFromModel(model *PositionModel) Position {
	return SolidPosition{
		model: model,
	}
}

func NewPosition(lat int, lon int) Position {
	return PositionFromModel(&PositionModel{lat, lon})
}

// Model

type PositionModel struct {
	Lat int `json:"lat"`
	Lon int `json:"lon"`
}

func EmptyPositionModel() *PositionModel {
	return &PositionModel{}
}

func (model *PositionModel) Change(newModel *PositionModel) {
	model.Lat = newModel.Lat
	model.Lon = newModel.Lon
}

func (model *PositionModel) LatF() float64 {
	return model.f(model.Lat)
}

func (model *PositionModel) LonF() float64 {
	return model.f(model.Lon)
}

func (model *PositionModel) f(v int) float64 {
	return float64(v) / 1000000
}

// Solid

type SolidPosition struct {
	model    *PositionModel
	position Position
}

func NewSolidPosition(model *PositionModel, Position Position) Position {
	return &SolidPosition{model, Position}
}

func (addr SolidPosition) Update(newModel *PositionModel) error {
	addr.model.Change(newModel)
	if addr.position == nil {
		return nil
	}
	return addr.position.Update(newModel)
}

func (addr SolidPosition) Model() *PositionModel {
	return addr.model
}
