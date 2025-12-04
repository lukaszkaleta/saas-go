package messages

import "github.com/lukaszkaleta/saas-go/universal"

type Message interface {
	Model() *Model
}

type Model struct {
	Id      int64                  `json:"id"`
	OwnerId int64                  `json:"owner_id"`
	Value   string                 `json:"value"`
	Actions universal.ActionsModel `json:"actions"`
}

func EmptyModel() *Model {
	return &Model{
		Id:      0,
		OwnerId: 0,
		Value:   "",
	}
}
