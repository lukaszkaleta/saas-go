package messages

import (
	"context"
)

type Inbox interface {
	Last(ctx context.Context) ([]Message, error)
	CountUnread(ctx context.Context) (int, error)
}
