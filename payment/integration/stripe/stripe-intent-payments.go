package stripe

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lukaszkaleta/saas-go/payment"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
)

type StripeIntentPayments struct {
	payments payment.Payments
}

func NewStripeIntentPayments(payments payment.Payments) payment.Payments {
	return &StripeIntentPayments{payments: payments}
}

func (s StripeIntentPayments) Create(ctx context.Context, id int64) (payment.Intent, error) {
	intent, err := s.createInternalIntent(ctx, id)
	if err != nil {
		return nil, err
	}

	stripeId, clientSecret, err := s.createStripePaymentIntent(ctx, intent)
	if err != nil {
		return nil, err
	}

	return s.enrichIntent(ctx, intent, stripeId, clientSecret)
}

func (s StripeIntentPayments) Search() payment.Search {
	return s.payments.Search()
}

func (s StripeIntentPayments) createInternalIntent(ctx context.Context, id int64) (payment.Intent, error) {
	return s.payments.Create(ctx, id)
}

func (s StripeIntentPayments) createStripePaymentIntent(ctx context.Context, internalIntent payment.Intent) (stripeID string, clientSecret string, err error) {

	internalIntentModel, err := internalIntent.Model(ctx)
	if err != nil {
		return
	}
	amount := internalIntentModel.Amount
	currency := internalIntentModel.Currency
	ref := internalIntentModel.Reference

	sAmount := stripe.Int64(amount)
	sCurrency := stripe.String(strings.ToLower(currency))
	slog.Info("Payment to stripe: ", "amount", sAmount, "currency", sCurrency)
	params := &stripe.PaymentIntentParams{
		Amount:   sAmount,
		Currency: sCurrency,

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

	slog.Info("Payment to stripe created: ", "client secret", pi.ClientSecret)
	return pi.ID, pi.ClientSecret, nil
}

func (s StripeIntentPayments) enrichIntent(
	ctx context.Context,
	intent payment.Intent,
	paymentId string,
	clientSecret string,
) (payment.Intent, error) {
	err := intent.UpdateStripeIntent(ctx, paymentId, clientSecret)
	if err != nil {
		return nil, err
	}
	return intent, nil
}
