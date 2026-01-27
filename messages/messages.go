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

func OwnerIds(ctx context.Context, list []Message) []int64 {
	ids := make([]int64, len(list))
	for _, msg := range list {
		ids = append(ids, msg.Model(ctx).OwnerId)
	}
	return ids
}

func Models(ctx context.Context, list []Message) []*MessageModel {
	models := make([]*MessageModel, len(list))
	for _, msg := range list {
		models = append(models, msg.Model(ctx))
	}
	return models
}

func ModelsAndOwners(ctx context.Context, list []Message) ([]*MessageModel, []int64) {
	models := make([]*MessageModel, len(list))
	ownerIds := make([]int64, len(list))
	for _, msg := range list {
		model := msg.Model(ctx)
		models = append(models, model)
		ownerIds = append(ownerIds, model.OwnerId)
	}
	return models, ownerIds
}
