package job

import "github.com/lukaszkaleta/saas-go/universal"

type Offers interface {
	Waiting() []*Offer

	Make(person universal.Person, model OfferModel)
}

type NoOffers struct {
}

func (n NoOffers) Waiting() []*Offer {
	return nil
}

func (n NoOffers) Make(person universal.Person, model OfferModel) {
}
