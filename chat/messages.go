package chat

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Messages interface {
	universal.Creator[string, Message]
	AddGenerated(ctx context.Context, value string) (Message, error)
	List(ctx context.Context) ([]Message, error)
	ById(ctx context.Context, id int64) (Message, error)
	Acknowledge(ctx context.Context) error
}
