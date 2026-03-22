package pg

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/payment"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgIntent struct {
	db        *pg.PgDb
	id        int64
	reference string
}

func (p PgIntent) ID() int64 {
	return p.id
}

func (p *PgIntent) Model(ctx context.Context) (*payment.IntentModel, error) {
	query := `
		SELECT 
			id, reference, stripe_payment_intent_id, stripe_client_secret, job_id, payer_id, payee_id, amount, currency, status, action_created_by_id, action_created_at
		FROM pay_payment_intent 
		WHERE reference = @reference
	`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"reference": p.reference})
	if err != nil {
		return nil, err
	}
	model, err := pgx.CollectOneRow(rows, MapIntentModel)
	if err != nil {
		return nil, err
	}
	p.id = model.Id
	return model, nil
}

func MapIntent(db *pg.PgDb) pgx.RowToFunc[payment.Intent] {
	return func(row pgx.CollectableRow) (payment.Intent, error) {
		model, err := MapIntentModel(row)
		if err != nil {
			return nil, err
		}
		return payment.NewSolidIntent(model, &PgIntent{db: db, id: model.Id, reference: model.Reference}, model.Id), nil
	}
}

func MapIntentModel(row pgx.CollectableRow) (*payment.IntentModel, error) {
	model := payment.EmptyIntentModel()

	stripePaymentIntentId := sql.NullString{}
	stripeClientSecret := sql.NullString{}

	actionCreatedModel := universal.EmptyCreatedActionModel()
	actions := make(map[string]*universal.ActionModel)
	actions[actionCreatedModel.Name] = actionCreatedModel

	err := row.Scan(
		&model.Id,
		&model.Reference,
		&stripePaymentIntentId,
		&stripeClientSecret,
		&model.JobId,
		&model.PayerId,
		&model.PayeeId,
		&model.Amount,
		&model.Currency,
		&model.Status,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
	)
	if err != nil {
		return nil, err
	}

	model.StripePaymentIntentId = stripePaymentIntentId.String
	model.StripeClientSecret = stripeClientSecret.String
	model.Actions = &universal.ActionsModel{List: actions}

	return model, nil
}
