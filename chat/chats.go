package chat

import (
	"context"
)

type Chats interface {
	Create(ctx context.Context, workerId int64) (Chat, error)
	Delete(ctx context.Context) error
	ById(ctx context.Context, id int64) (Chat, error)
	LastMessages(ctx context.Context) ([]Message, error)
	ByWorkerId(ctx context.Context, id int64) (Chat, error)
}
