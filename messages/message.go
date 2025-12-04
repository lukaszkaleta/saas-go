package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Message interface {
	Model(ctx context.Context) *Model
}

type Model struct {
	Id      int64                  `json:"id"`
	OwnerId int64                  `json:"owner_id"`
	Value   string                 `json:"value"`
	Actions universal.ActionsModel `json:"actions"`
}

func EmptyModel() *Model {
	return EmptyOwnerModel(0)
}

func EmptyOwnerModel(ownerId int64) *Model {
	return &Model{
		Id:      0,
		OwnerId: ownerId,
		Value:   "",
	}
}
