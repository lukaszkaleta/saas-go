package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type OLDMessages interface {
	universal.Lister[OLDMessage]
	universal.Idables[OLDMessage]
	// Adder as default implementation, if you want to intercept
	universal.Adder[*OLDMessageModel, OLDMessage]
	AddSimple(ctx context.Context, recipientId int64, value string) (OLDMessage, error)
	AddGenerated(ctx context.Context, recipientId int64, value string) (OLDMessage, error)
	ForRecipient(ctx context.Context, recipient universal.Idable) ([]OLDMessage, error)
	ForRecipientById(ctx context.Context, id int64) ([]OLDMessage, error)
	Acknowledge(ctx context.Context) error
	LastQuestions(ctx context.Context) ([]OLDMessage, error)
	Delete(ctx context.Context) error
}

func OwnerIds(ctx context.Context, list []OLDMessage) []int64 {
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

func OLDInvolvedUserIds(ctx context.Context, list []OLDMessage) []*int64 {
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

func OLDModels(ctx context.Context, list []OLDMessage) []*OLDMessageModel {
	models := make([]*OLDMessageModel, len(list))
	for i, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		models[i] = model
	}
	return models
}

func OLDModelsAndOwners(ctx context.Context, list []OLDMessage) ([]*OLDMessageModel, []int64) {
	models := make([]*OLDMessageModel, len(list))
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
