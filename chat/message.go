package chat

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Message interface {
	universal.Idable
	Model(ctx context.Context) (*MessageModel, error)
}

type MessageModel struct {
	Id             int64                   `json:"id"`
	ChatId         int64                   `json:"chatId"`
	Value          string                  `json:"value"`
	ValueGenerated bool                    `json:"valueGenerated"`
	Actions        *universal.ActionsModel `json:"actions"`
}

func (m MessageModel) ID() int64 {
	return m.Id
}

type SolidMessage struct {
	model   *MessageModel
	message Message
}

func NewSolidMessage(model *MessageModel, message Message) Message {
	return &SolidMessage{
		model:   model,
		message: message,
	}
}

func (m *SolidMessage) ID() int64 {
	return m.message.ID()
}

func (m *SolidMessage) Model(ctx context.Context) (*MessageModel, error) {
	return m.model, nil
}
