package pgjob

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgOffers struct {
	db    *pg.PgDb
	JobId int64
}

func (p *PgOffers) Waiting(ctx context.Context) ([]job.Offer, error) {
	query := "select * from job_offer where job_id = $1 and action_accepted_at is null and action_rejected_at is null"
	rows, err := p.db.Pool.Query(ctx, query, p.JobId)
	if err != nil {
		return nil, err
	}
	return MapOffers(rows, p.db)
}

func (pgOffers *PgOffers) Make(ctx context.Context, model *job.OfferModel) (job.Offer, error) {
	offerId := int64(0)
	user := user.CurrentUser(ctx)

	query := "INSERT INTO job_offer (job_id, price_value, price_currency, description_value, action_created_by_id) VALUES( $1, $2, $3, $4, $5 ) returning id"
	row := pgOffers.db.Pool.QueryRow(
		ctx,
		query,
		pgOffers.JobId,
		model.Price.Value,
		model.Price.Currency,
		model.Description.Value,
		user.Id,
	)
	err := row.Scan(&offerId)
	if err != nil {
		return nil, err
	}
	pgOffer := PgOffer{
		db: pgOffers.db,
		Id: offerId,
	}
	actionsList := make(map[string]*universal.ActionModel)
	actionsList[job.Created] = &universal.ActionModel{
		ById:   &user.Id,
		MadeAt: time.Now(),
		Name:   job.Created,
	}
	return job.NewSolidOffer(
		&job.OfferModel{
			Id:          offerId,
			Description: model.Description,
			Price:       model.Price,
			Actions:     universal.ActionsModel{List: actionsList},
		},
		&pgOffer,
	), nil
}

func MapOffers(rows pgx.Rows, db *pg.PgDb) ([]job.Offer, error) {
	offers := []job.Offer{}
	id := int64(0)
	for rows.Next() {
		pgOffer := &PgOffer{db: db, Id: id}
		offerModel, err := MapOffer(rows)
		if err != nil {
			return nil, err
		}
		solidOffer := job.NewSolidOffer(offerModel, pgOffer)
		offers = append(offers, solidOffer)
	}
	return offers, nil
}
