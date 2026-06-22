package chat

import (
	"context"

	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Messages interface {
	universal.Creator[string, Message]
	AddGenerated(ctx context.Context, value string, reason string) (Message, error)
	List(ctx context.Context) ([]Message, error)
	ById(ctx context.Context, id int64) (Message, error)
	Acknowledge(ctx context.Context) error
	LastReadMessageId(ctx context.Context) (int64, error)
	LastReadMessageAt(ctx context.Context) (time.Time, error)
}

func AllMessages(ctx context.Context, messages Messages) ([]Message, []int64, error) {
	list, err := messages.List(ctx)
	if err != nil {
		return nil, nil, err
	}
	at, err := messages.LastReadMessageAt(ctx)
	if err != nil {
		return nil, nil, err
	}

	currentUserId := universal.CurrentUserId(ctx)
	var newMessagesIds []int64
	for _, m := range list {
		model, err := m.Model(ctx)
		if err != nil {
			return nil, nil, err
		}
		if model.Actions.Created().MadeAt.After(at) {
			if currentUserId != nil && model.Actions.Created().ById != nil && *currentUserId == *model.Actions.Created().ById {
				continue
			}
			newMessagesIds = append(newMessagesIds, model.Id)
		}
	}
	return list, newMessagesIds, nil
}
