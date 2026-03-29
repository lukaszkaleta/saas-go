package stripe

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/lukaszkaleta/saas-go/payment"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type StripeSessionPayments struct {
	payments   payment.Payments
	successUrl string
	failureUrl string
}

func NewStripeSessionPayments(payments payment.Payments, successUrl string, failureUrl string) payment.Payments {
	return &StripeSessionPayments{payments: payments, successUrl: successUrl, failureUrl: failureUrl}
}

func (s *StripeSessionPayments) Create(ctx context.Context, id int64) (payment.Intent, error) {
	intent, err := s.createInternalSession(ctx, id)
	if err != nil {
		return nil, err
	}

	stripeId, stripeUrl, err := s.createStripeSession(ctx, intent)
	if err != nil {
		return nil, err
	}

	return s.enrichIntent(ctx, intent, stripeId, stripeUrl)
}

func (s *StripeSessionPayments) Search() payment.Search {
	return s.payments.Search()
}

func (s *StripeSessionPayments) createInternalSession(ctx context.Context, id int64) (payment.Intent, error) {
	return s.payments.Create(ctx, id)
}

func (s *StripeSessionPayments) createStripeSession(ctx context.Context, internalIntent payment.Intent) (stripeId string, stripeUrl string, err error) {

	internalIntentModel, err := internalIntent.Model(ctx)
	if err != nil {
		return
	}
	amount := internalIntentModel.Amount
	currency := internalIntentModel.Currency
	ref := internalIntentModel.Reference
	jobId := internalIntentModel.JobId

	sAmount := stripe.Int64(amount)
	sCurrency := stripe.String(strings.ToLower(currency))
	slog.Info("Payment to stripe: ", "amount", sAmount, "currency", sCurrency)

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),

		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"mobilepay",
			"vipps", // enable if available
		}),

		Metadata: map[string]string{
			"payment_id": ref,
			"job_id":     strconv.FormatInt(jobId, 10),
		},

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   sCurrency,
					UnitAmount: sAmount,
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Naborly Job Payment"),
					},
				},
			},
		},

		SuccessURL: stripe.String(s.successUrl),
		CancelURL:  stripe.String(s.failureUrl),
	}

	//params.AddHeader("Vipps-Preview", "v1")
	params.Params.Headers = http.Header{
		"Vipps-Preview": []string{"v1"},
	}

	stripeSession, err := session.New(params)
	if err != nil {
		log.Println("stripe error:", err)
		return "", "", err
	}

	slog.Info("Session in stripe created: ", "url", stripeSession.URL)
	return stripeSession.ID, stripeSession.URL, nil
}

func (s *StripeSessionPayments) enrichIntent(
	ctx context.Context,
	intent payment.Intent,
	paymentId string,
	url string,
) (payment.Intent, error) {
	err := intent.UpdateStripeSession(ctx, paymentId, url)
	if err != nil {
		return nil, err
	}
	return intent, nil
}
