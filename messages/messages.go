package messages

import "context"

type Messages interface {
	Add(ctx context.Context, model *Model) (Message, error)
}
