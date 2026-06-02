package finance

import (
	"time"
)

type LedgerEntry struct {
	ID                    int64          `json:"id"`
	SellerID              int64          `json:"seller_id"`
	BuyerID               *int64         `json:"buyer_id,omitzero"`
	JobID                 *int64         `json:"job_id,omitzero"`
	Type                  EventType      `json:"type"`
	Amount                int64          `json:"amount"`
	GrossAmount           *int64         `json:"gross_amount,omitzero"`
	FeeAmount             *int64         `json:"fee_amount,omitzero"`
	NetAmount             *int64         `json:"net_amount,omitzero"`
	Currency              string         `json:"currency"`
	StripePaymentIntentID *string        `json:"stripe_payment_intent_id,omitzero"`
	StripeTransferID      *string        `json:"stripe_transfer_id,omitzero"`
	StripePayoutID        *string        `json:"stripe_payout_id,omitzero"`
	StripeRefundID        *string        `json:"stripe_refund_id,omitzero"`
	OccurredAt            time.Time      `json:"occurred_at"`
	CreatedAt             time.Time      `json:"created_at"`
	Metadata              map[string]any `json:"metadata,omitzero"`
}
