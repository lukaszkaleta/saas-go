package messages

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
	OwnerId        int64                   `json:"ownerId"`
	RecipientId    int64                   `json:"recipientId"`
	Value          string                  `json:"value"`
	ValueGenerated bool                    `json:"generated"`
	Actions        *universal.ActionsModel `json:"actions"`
}

func EmptyModel() *MessageModel {
	return EmptyOwnerModel(0)
}

func (m MessageModel) ID() int64 {
	return m.Id
}

func (m MessageModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

func EmptyOwnerModel(ownerId int64) *MessageModel {
	return &MessageModel{
		Id:          0,
		OwnerId:     ownerId,
		RecipientId: 0,
		Value:       "",
		Actions:     universal.EmptyActionsModel(),
	}
}

// Solid

type SolidMessage struct {
	Id      int64
	model   *MessageModel
	message Message
}

func (m *SolidMessage) Acknowledge(ctx context.Context) error {
	if m.message != nil {
		return m.message.Acknowledge(ctx)
	}
	return nil
}

func NewSolidMessage(model *MessageModel, message Message, id int64) Message {
	return &SolidMessage{
		Id:      id,
		model:   model,
		message: message,
	}
}

func (m *SolidMessage) Model(ctx context.Context) (*MessageModel, error) {
	return m.model, nil
}

func (m *SolidMessage) ID() int64 {
	return m.Id
}
