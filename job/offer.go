package job

import "github.com/lukaszkaleta/saas-go/universal"

type Offer interface {
	Accept(person universal.Person) error
	Reject(person universal.Person) error
}

type OfferModel struct {
}
