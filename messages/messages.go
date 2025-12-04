package messages

import "context"

type Messages interface {
	Add(ctx context.Context, value string) (Message, error)
	AddFromModel(ctx context.Context, model *Model) (Message, error)
}
