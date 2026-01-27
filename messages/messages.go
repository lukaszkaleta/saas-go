package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Messages interface {
	universal.Lister[Message]
	Add(ctx context.Context, recipientId int64, value string) (Message, error)
	AddFromModel(ctx context.Context, model *MessageModel) (Message, error)
}
