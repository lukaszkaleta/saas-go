package universal

import (
	"context"
)

// Generic inbox, something where you can read last
type Inbox[T Idable] interface {
	Last(ctx context.Context) ([]T, error)
	CountUnread(ctx context.Context) (int, error)
}
