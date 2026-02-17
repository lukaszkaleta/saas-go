package universal

import (
	"context"
	"time"
)

// Objects which holds actions are actions aware

type ActionsAware interface {
	Actions() Actions
}

type ActionsAwareModel interface {
	Idable
	GetActions() *ActionsModel
}

func (am *ActionsModel) GetActions() *ActionsModel {
	return am
}

type Actions interface {
	List() map[string]*Action
	WithName(name string) Action
	Created() Action
	Model(ctx context.Context) (*ActionsModel, error)
}

type ActionsModel struct {
	Idable
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

func (am *ActionsModel) Created() *ActionModel {
	return am.WithName("created")
}

func (am *ActionsModel) WithName(name string) *ActionModel {
	return am.List[name]
}

func (am *ActionsModel) ID() int64 {
	return 0
}

func (am *ActionsModel) CreatedById() *int64 {
	return am.Created().ById
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

func CreatedById[T ActionsAwareModel](ctx context.Context, instance any) (int64, error) {
	_, ok := instance.(ActionsAware)
	if !ok {
		return 0, nil
	}
	modelAware, ok := instance.(ModelAware[T])
	if !ok {
		return 0, nil
	}
	model, err2 := modelAware.Model(ctx)
	if err2 != nil {
		return 0, err2
	}
	actionsAwareModel, ok := any(model).(ActionsAwareModel)
	if !ok {
		return 0, nil
	}
	userId := actionsAwareModel.GetActions().CreatedById()
	if userId == nil {
		return 0, nil
	}
	return *userId, nil
}

func CreatedByIdFromModel(model ActionsAwareModel) int64 {
	userId := model.GetActions().CreatedById()
	if userId == nil {
		return 0
	}
	return *userId
}

func CreatedByIdFromModels(models []ActionsAwareModel) ([]int64, error) {
	ids := make([]int64, len(models))
	for i, model := range models {
		id := CreatedByIdFromModel(model)
		ids[i] = id
	}
	return ids, nil
}
