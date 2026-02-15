package universal

import (
	"context"
	"time"
)

type Action interface {
	Model() *ActionModel
}

type ActionModel struct {
	ById   *int64     `json:"byId"`
	MadeAt *time.Time `json:"at"`
	Name   string     `json:"name"`
}

func (m *ActionModel) Exists() bool {
	return m.ById != nil && *m.ById > 0 && !m.MadeAt.IsZero()
}

func EmptyActionModel(name string) *ActionModel {
	noOne := int64(0)
	return NowActionModelForUser(name, &noOne)
}

func ZeroActionModelForUser(name string, userId *int64) *ActionModel {
	return &ActionModel{
		ById:   userId,
		MadeAt: &time.Time{},
		Name:   name,
	}
}

func NowActionModelForUser(name string, userId *int64) *ActionModel {
	now := time.Now()
	return &ActionModel{
		ById:   userId,
		MadeAt: &now,
		Name:   name,
	}
}

func EmptyCreatedActionModel() *ActionModel {
	return EmptyActionModel("created")
}

type SolidAction struct {
	model  *ActionModel
	action Action
}

func (action *SolidAction) Model() *ActionModel {
	return action.model
}

func CurrentUserId(ctx context.Context) *int64 {
	return ctx.Value("current-user-id").(*int64)
}
