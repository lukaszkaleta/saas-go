package chat

import (
	"context"
)

type Messages interface {
	List(ctx context.Context) ([]Message, error)
}
