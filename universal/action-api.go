package universal

import (
	"context"
	"time"
)

type Action interface {
	Model() ActionsModel
}

type ActionModel struct {
	ById   *int64
	MadeAt time.Time
	Name   string
}

func (m *ActionModel) Exists() bool {
	return m.ById != nil && *m.ById > 0 && !m.MadeAt.IsZero()
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
