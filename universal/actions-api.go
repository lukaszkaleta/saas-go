package universal

import (
	"context"
	"time"
)

type Actions interface {
	List() map[string]*Action
	WithName(name string) Action
	Created() Action
	Model() *ActionsModel
}

type ActionsModel struct {
	List map[string]*ActionModel `json:"list"`
}

func EmptyActionsModel() *ActionsModel {
	return &ActionsModel{List: make(map[string]*ActionModel)}
}

func CreatedNowActions(ctx context.Context) *ActionsModel {
	now := time.Now()
	model := EmptyActionsModel()
	createdModel := EmptyCreatedActionModel()
	createdModel.MadeAt = &now
	createdModel.ById = CurrentUserId(ctx)
	model.List["created"] = createdModel
	return model
}

func (am ActionsModel) Created() *ActionModel {
	return am.WithName("created")
}

func (am ActionsModel) WithName(name string) *ActionModel {
	return am.List[name]
}

type SolidActions struct {
	actions Actions
	model   *ActionsModel
}

func (s *SolidActions) List() map[string]*Action {
	return s.actions.List()
}

func (s *SolidActions) Model() *ActionsModel {
	return s.model
}

func (s *SolidActions) WithName(name string) Action {
	return s.actions.WithName(name)
}
func (s *SolidActions) Created() Action {
	return s.WithName("created")
}
