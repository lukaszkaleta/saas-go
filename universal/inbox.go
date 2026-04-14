package universal

import (
	"context"
)

// Generic inbox, something where you can read last
type Inbox[T Idable] interface {
	Last(ctx context.Context) ([]T, error)
	CountUnread(ctx context.Context) (int, error)
}

// Generic outbox, something where you can read what you have sent
type Outbox[T Idable] interface {
	Last(ctx context.Context) ([]T, error)
	Old(ctx context.Context) ([]T, error)
	CountUnread(ctx context.Context) (int, error)
}
