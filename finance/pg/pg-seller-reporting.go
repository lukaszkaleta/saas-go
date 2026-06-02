package pg

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/finance"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgSellerReporting struct {
	db       *pg.PgDb
	sellerID int64
}

func (r *PgSellerReporting) SumInPeriod(ctx context.Context, period universal.DateRange) (int64, error) {
	const query = `
		SELECT SUM(amount)
		FROM financial_ledger
		WHERE seller_id = $1
		  AND type = $2
		  AND occurred_at BETWEEN $3 AND $4
	`

	var sum *int64
	err := r.db.Pool.QueryRow(ctx, query, r.sellerID, finance.EventPayoutRelease, period.From, period.To).Scan(&sum)
	if err != nil {
		return 0, err
	}

	if sum == nil {
		return 0, nil
	}

	return *sum, nil
}
