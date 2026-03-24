package pg

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/payment"
)

var ErrInvalidAmount = errors.New("invalid amount")
var ErrNoAcceptedOffer = errors.New("missing accepted offer")

type PgPayments struct {
	db  *pg.PgDb
	job job.Job
}

func NewPgPayments(db *pg.PgDb, job job.Job) payment.Payments {
	return &PgPayments{db: db, job: job}
}

func (p PgPayments) Create(ctx context.Context, offerId int64) (payment.Intent, error) {

	ref := uuid.NewString()
	accepted, err := p.job.Offers().ById(ctx, offerId)
	if err != nil {
		return nil, err
	}
	if accepted == nil {
		return nil, ErrNoAcceptedOffer
	}
	offerModel, err := accepted.Model(ctx)
	if err != nil {
		return nil, err
	}
	amount := offerModel.Price.Value
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	payeeId := offerModel.Actions.CreatedById()
	model, err := p.job.Model(ctx)
	if err != nil {
		return nil, err
	}
	payerId := model.Actions.CreatedById()

	const query = `
		INSERT INTO pay_payment_intent (
			reference,
			job_id,
			payer_id,
			payee_id,
			amount,
			currency,
			status,
			action_created_by_id,
			action_created_at
		)
		VALUES (
			@reference,
			@job_id,
			@payer_id,
			@payee_id,
			@amount,
			@currency,
			@status,
			@created_by,
			now()
		)
		RETURNING
			id
	`

	var id int64
	err = p.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"reference":  ref,
		"job_id":     p.job.ID(),
		"payer_id":   payerId,
		"payee_id":   payeeId,
		"amount":     amount,
		"currency":   "NOK",
		"status":     "CREATED",
		"created_by": payerId,
	}).Scan(&id)

	if err != nil {
		return nil, err
	}

	return payment.NewSolidIntent(
		&payment.IntentModel{
			Id:        id,
			Reference: ref,
			JobId:     p.job.ID(),
			PayerId:   *payerId,
			PayeeId:   *payeeId,
			Amount:    int64(amount),
			Currency:  "NOK",
			Status:    "CREATED",
		},
		&PgIntent{db: p.db, id: id, reference: ref},
		id,
	), nil
}

func (p PgPayments) Search() payment.Search {
	return NewPgSearch(p.db)
}
