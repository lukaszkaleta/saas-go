package chat

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Message interface {
	universal.Idable
	Model(ctx context.Context) (*MessageModel, error)
	Acknowledge(ctx context.Context) error
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

func (m MessageModel) GetActions() *universal.ActionsModel {
	return m.Actions
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

func (m *SolidMessage) Acknowledge(ctx context.Context) error {
	return m.message.Acknowledge(ctx)
}
