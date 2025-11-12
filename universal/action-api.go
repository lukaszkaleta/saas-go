package universal

import "time"

type Action interface {
	Model() ActionsModel
}

type ActionModel struct {
	ById   int64
	MadeAt time.Time
	Name   string
}

type SolidAction struct {
	model  *ActionModel
	action Action
}

func (action *SolidAction) Model() *ActionModel {
	return action.model
}
