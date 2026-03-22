package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/payment"
)

type PgSearch struct {
	db *pg.PgDb
}

func NewPgSearch(db *pg.PgDb) payment.Search {
	return &PgSearch{db: db}
}

func (p PgSearch) Intent(ctx context.Context, reference string) (payment.Intent, error) {
	query := `
		SELECT 
			id, reference, stripe_payment_intent_id, stripe_client_secret, job_id, payer_id, payee_id, amount, currency, status, action_created_by_id, action_created_at
		FROM pay_payment_intent 
		WHERE reference = @reference
	`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"reference": reference})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapIntent(p.db))
}
