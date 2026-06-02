package pg

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/finance"
)

type PgFinancialLedger struct {
	db *pg.PgDb
}

func NewPgFinancialLedger(db *pg.PgDb) finance.FinancialLedger {
	return &PgFinancialLedger{db: db}
}

func (l *PgFinancialLedger) JobView(ctx context.Context, jobID int64) ([]finance.LedgerEntry, error) {
	const query = `
		SELECT
			id,
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
			created_at,
			metadata
		FROM financial_ledger
		WHERE job_id = $1
		ORDER BY occurred_at
	`

	rows, err := l.db.Pool.Query(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []finance.LedgerEntry
	for rows.Next() {
		var entry finance.LedgerEntry
		err := rows.Scan(
			&entry.ID,
			&entry.SellerID,
			&entry.BuyerID,
			&entry.JobID,
			&entry.Type,
			&entry.Amount,
			&entry.GrossAmount,
			&entry.FeeAmount,
			&entry.NetAmount,
			&entry.Currency,
			&entry.StripePaymentIntentID,
			&entry.StripeTransferID,
			&entry.StripePayoutID,
			&entry.StripeRefundID,
			&entry.OccurredAt,
			&entry.CreatedAt,
			&entry.Metadata,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (l *PgFinancialLedger) SellerReporting(sellerID int64) finance.SellerReporting {
	return &PgSellerReporting{
		db:       l.db,
		sellerID: sellerID,
	}
}
