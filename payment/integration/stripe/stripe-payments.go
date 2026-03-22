package stripe

import (
	"context"

	"github.com/lukaszkaleta/saas-go/payment"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
)

type StripePayments struct {
	payments payment.Payments
}

func (s StripePayments) Create(ctx context.Context) (payment.Intent, error) {
	intent, err := s.createInternalIntent(ctx)
	if err != nil {
		return nil, err
	}

	stripeId, clientSecret, err := s.createStripePaymentIntent(ctx, intent)
	if err != nil {
		return nil, err
	}

	return s.enrichIntent(ctx, intent, stripeId, clientSecret)
}

func (s StripePayments) Search() payment.Search {
	//TODO implement me
	panic("implement me")
}

func (s StripePayments) createInternalIntent(ctx context.Context) (payment.Intent, error) {
	return s.payments.Create(ctx)
}

func (s StripePayments) createStripePaymentIntent(ctx context.Context, internalIntent payment.Intent) (stripeID string, clientSecret string, err error) {

	internalIntentModel, err := internalIntent.Model(ctx)
	if err != nil {
		return
	}
	amount := internalIntentModel.Amount
	currency := internalIntentModel.Currency
	ref := internalIntentModel.Reference

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),

		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},

		Metadata: map[string]string{
			"internal_intent_id": ref,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", err
	}

	return pi.ID, pi.ClientSecret, nil
}

func (s StripePayments) enrichIntent(
	ctx context.Context,
	intent payment.Intent,
	paymentId string,
	clientSecret string,
) (payment.Intent, error) {
	err := intent.UpdateStripe(ctx, paymentId, clientSecret)
	if err != nil {
		return nil, err
	}
	return intent, nil
}
