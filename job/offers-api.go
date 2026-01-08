package job

import (
	"context"
)

type Offers interface {
	Waiting(ctx context.Context) ([]Offer, error)
	Make(ctx context.Context, model *OfferModel) (Offer, error)
}

type NoOffers struct {
}

func (n NoOffers) Waiting(ctx context.Context) ([]Offer, error) {
	return nil, nil
}

func (n NoOffers) Make(ctx context.Context, model *OfferModel) (Offer, error) {
	return nil, nil
}
