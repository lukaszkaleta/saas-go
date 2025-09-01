package universal

// API

type Radar interface {
	Model() *RadarModel
	Update(newModel *RadarModel) error
}

// Builder

func RadarFromModel(model *RadarModel) Radar {
	return SolidRadar{
		model: model,
	}
}

// Model

type RadarModel struct {
	Position  *PositionModel `json:"position"`
	Perimeter int            `json:"perimeter"`
}

func NewRadarModel() *RadarModel {
	return &RadarModel{
		Position: &PositionModel{},
	}
}

func EmptyRadarModel() *RadarModel {
	return &RadarModel{
		Position:  EmptyPositionModel(),
		Perimeter: 0,
	}
}

func (model *RadarModel) Change(newModel *RadarModel) {
	model.Position.Change(newModel.Position)
	model.Perimeter = newModel.Perimeter
}

// Solid

type SolidRadar struct {
	model *RadarModel
	radar Radar
}

func NewSolidRadar(model *RadarModel, Radar Radar) Radar {
	return &SolidRadar{model, Radar}
}

func (radar SolidRadar) Update(newModel *RadarModel) error {
	radar.model.Change(newModel)
	if radar.radar == nil {
		return nil
	}
	return radar.radar.Update(newModel)
}

func (radar SolidRadar) Model() *RadarModel {
	return radar.model
}
