package stripe

import (
	"context"

	"github.com/lukaszkaleta/saas-go/payment"
)

type StripePayments struct {
	payments payment.Payments
}

func (s StripePayments) Create(ctx context.Context) (payment.Intent, error) {
	create, err := s.payments.Create(ctx)
	if err != nil {
		return nil, err
	}
	// call stripe.
	return create, nil
}

func (s StripePayments) Search() payment.Search {
	//TODO implement me
	panic("implement me")
}
