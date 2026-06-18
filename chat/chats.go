package chat

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Chats interface {
	universal.Deleter
	Ensure(ctx context.Context, workerId int64) (Chat, error)
	ById(ctx context.Context, id int64) (Chat, error)
	LastMessages(ctx context.Context) ([]Message, error)
	ByWorkerId(ctx context.Context, id int64) (Chat, error)
}
