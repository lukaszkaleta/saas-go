package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Message interface {
	universal.Idable
	Model(ctx context.Context) *MessageModel
	Acknowledge(ctx context.Context) error
}

type MessageModel struct {
	Id          int64                   `json:"id"`
	OwnerId     int64                   `json:"ownerId"`
	RecipientId int64                   `json:"recipientId"`
	Value       string                  `json:"value"`
	Actions     *universal.ActionsModel `json:"actions"`
}

func EmptyModel() *MessageModel {
	return EmptyOwnerModel(0)
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

func (m *SolidMessage) Model(ctx context.Context) *MessageModel {
	return m.model
}

func (m *SolidMessage) ID() int64 {
	return m.Id
}
