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
}

// No offers implementation

type NoOffers struct {
}

func (n NoOffers) Waiting(ctx context.Context) ([]Offer, error) {
	return nil, nil
}

func (n NoOffers) Make(ctx context.Context, model *OfferModel) (Offer, error) {
	return nil, nil
}

// Offer maker with

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
	offerModel := offer.Model()
	offerMessage := offerModel.Description.Value
	// Make offer message a message which will be put into chat:
	message := fmt.Sprintf("% %s, %s",
		offerModel.Price.Value,
		offerModel.Price.Currency,
		offerMessage,
	)
	userCreatedAt, err := universal.CreatedUserId(ctx, m.job)
	if err != nil {
		return nil, err
	}
	_, err = m.job.Messages().Add(ctx, *userCreatedAt, message)
	if err != nil {
		return nil, err
	}
	return offer, nil
}
