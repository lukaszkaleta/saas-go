package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/finance"
)

type PgFinancialLedger struct {
	db *pg.PgDb
}

func NewPgFinancialLedger(db *pg.PgDb) finance.FinancialLedger {
	return &PgFinancialLedger{db: db}
}

func (l *PgFinancialLedger) Record(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	const query = `
		INSERT INTO financial_ledger (
			seller_id,
			buyer_id,
			job_id,
			type,
			amount,
			gross_amount,
			fee_amount,
			net_amount,
			currency,
			stripe_payment_intent_id,
			stripe_transfer_id,
			stripe_payout_id,
			stripe_refund_id,
			occurred_at,
			metadata
		)
		VALUES (
			@seller_id,
			@buyer_id,
			@job_id,
			@type,
			@amount,
			@gross_amount,
			@fee_amount,
			@net_amount,
			@currency,
			@stripe_payment_intent_id,
			@stripe_transfer_id,
			@stripe_payout_id,
			@stripe_refund_id,
			@occurred_at,
			@metadata
		)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"seller_id":                entry.SellerID,
		"buyer_id":                 entry.BuyerID,
		"job_id":                   entry.JobID,
		"type":                     string(entry.Type),
		"amount":                   entry.Amount,
		"gross_amount":             entry.GrossAmount,
		"fee_amount":               entry.FeeAmount,
		"net_amount":               entry.NetAmount,
		"currency":                 entry.Currency,
		"stripe_payment_intent_id": entry.StripePaymentIntentID,
		"stripe_transfer_id":       entry.StripeTransferID,
		"stripe_payout_id":         entry.StripePayoutID,
		"stripe_refund_id":         entry.StripeRefundID,
		"occurred_at":              entry.OccurredAt,
		"metadata":                 entry.Metadata,
	}

	var id int64
	err := l.db.Pool.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (l *PgFinancialLedger) EscrowHold(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventEscrowHold
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) PayoutRelease(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventPayoutRelease
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) PlatformFee(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventPlatformFee
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) Payout(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventPayout
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) Refund(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventRefund
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) Chargeback(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventChargeback
	return l.Record(ctx, entry)
}

func (l *PgFinancialLedger) Adjustment(ctx context.Context, entry finance.LedgerEntry) (int64, error) {
	entry.Type = finance.EventAdjustment
	return l.Record(ctx, entry)
}
