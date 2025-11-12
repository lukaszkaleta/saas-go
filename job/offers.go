package job

import (
	"context"
)

type Offers interface {
	Waiting() []Offer
	Make(context context.Context, model OfferModel) (Offer, error)
}

type NoOffers struct {
}

func (n NoOffers) Waiting() []Offer {
	return nil
}

func (n NoOffers) Make(context context.Context, model OfferModel) (Offer, error) {
	return nil, nil
}
