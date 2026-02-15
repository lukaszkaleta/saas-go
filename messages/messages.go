package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Messages interface {
	universal.Lister[Message]
	universal.Idables[Message]
	Add(ctx context.Context, recipientId int64, value string) (Message, error)
	AddGenerated(ctx context.Context, recipientId int64, value string) (Message, error)
	AddFromModel(ctx context.Context, model *MessageModel) (Message, error)
	ForRecipient(ctx context.Context, recipient universal.Idable) ([]Message, error)
	ForRecipientById(ctx context.Context, id int64) ([]Message, error)
	Acknowledge(ctx context.Context) error
}

func OwnerIds(ctx context.Context, list []Message) []int64 {
	ids := make([]int64, len(list))
	for _, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		ids = append(ids, model.OwnerId)
	}
	return ids
}

func InvolvedUserIds(ctx context.Context, list []Message) []*int64 {
	idsMap := map[*int64]bool{}
	for _, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		id1 := model.RecipientId
		idsMap[&id1] = true
		id2 := model.Actions.Created().ById
		idsMap[id2] = true
	}
	ids := make([]*int64, 0, len(idsMap))
	for id := range idsMap {
		ids = append(ids, id)
	}
	return ids
}

func Models(ctx context.Context, list []Message) []*MessageModel {
	models := make([]*MessageModel, len(list))
	for i, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		models[i] = model
	}
	return models
}

func ModelsAndOwners(ctx context.Context, list []Message) ([]*MessageModel, []int64) {
	models := make([]*MessageModel, len(list))
	ownerIds := make([]int64, len(list))
	for i, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		models[i] = model
		ownerIds[i] = model.OwnerId
	}
	return models, ownerIds
}
