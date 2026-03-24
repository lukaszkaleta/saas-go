package job

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/universal"
)

type OfferMaker interface {
	Make(ctx context.Context, model *OfferModel) (Offer, error)
}

type OfferWaiter interface {
	Waiting(ctx context.Context) ([]Offer, error)
}

type Offers interface {
	OfferMaker
	OfferWaiter

	ById(ctx context.Context, id int64) (Offer, error)
	FromUser(ctx context.Context, user universal.Idable) (Offer, error)

	Accepted(ctx context.Context) (Offer, error)
}

// No offers implementation

type NoOffers struct {
}

func (n NoOffers) ById(ctx context.Context, id int64) (Offer, error) {
	return nil, nil
}

func (n NoOffers) FromUser(ctx context.Context, user universal.Idable) (Offer, error) {
	return nil, nil
}
func (n NoOffers) Waiting(ctx context.Context) ([]Offer, error) {
	return nil, nil
}

func (n NoOffers) Make(ctx context.Context, model *OfferModel) (Offer, error) {
	return nil, nil
}
func (n NoOffers) Accepted(ctx context.Context) (Offer, error) {
	return nil, nil
}

type MessagesOfferMaker struct {
	inner OfferMaker
	job   Job
}

func NewMessagesOfferMaker(inner OfferMaker, job Job) *MessagesOfferMaker {
	return &MessagesOfferMaker{
		inner: inner,
		job:   job,
	}
}

func (m MessagesOfferMaker) Make(ctx context.Context, model *OfferModel) (Offer, error) {
	offer, err := m.inner.Make(ctx, model)
	if err != nil {
		return nil, err
	}
	offerModel, err := offer.Model(ctx)
	if err != nil {
		return nil, err
	}
	offerMessage := offerModel.Description.Value
	// Make offer message a message which will be put into chat:
	message := fmt.Sprintf("%s: %s",
		offerModel.Price.UserFriendly(),
		offerMessage,
	)

	jobCreatedById, err := universal.CreatedById[OfferModel](ctx, offer)
	if err != nil {
		return nil, err
	}
	_, err = m.job.Messages().Add(ctx, jobCreatedById, message)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func ModelsAndOwners(ctx context.Context, list []Offer) ([]*OfferModel, []int64) {
	models := make([]*OfferModel, len(list))
	ownerIds := make([]int64, len(list))
	for i, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		models[i] = model
		id := model.Actions.CreatedById()
		ownerIds[i] = *id
	}
	return models, ownerIds
}
