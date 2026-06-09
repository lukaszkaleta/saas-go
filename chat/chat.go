package chat

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Chat interface {
	universal.Idable
	Model(ctx context.Context) (*ChatModel, error)
	Messages() Messages
}

type ChatModel struct {
	Id       int64                   `json:"id"`
	JobId    int64                   `json:"jobId"`
	WorkerId int64                   `json:"workerId"`
	Actions  *universal.ActionsModel `json:"actions"`
}

func (m ChatModel) ID() int64 {
	return m.Id
}

type SolidChat struct {
	Id       int64
	model    *ChatModel
	chat     Chat
	messages Messages
}

func NewSolidChat(model *ChatModel, chat Chat, id int64, messages Messages) Chat {
	return &SolidChat{
		Id:       id,
		model:    model,
		chat:     chat,
		messages: messages,
	}
}

func (c *SolidChat) ID() int64 {
	return c.Id
}

func (c *SolidChat) Model(ctx context.Context) (*ChatModel, error) {
	return c.model, nil
}

func (c *SolidChat) Messages() Messages {
	return c.messages
}
