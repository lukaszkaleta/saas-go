package chat

import (
	"context"
)

type ChatsApi interface {
	Create(ctx context.Context, workerId int64) (Chat, error)
}
