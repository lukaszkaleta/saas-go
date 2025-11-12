package pgjob

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgOffers struct {
	Db    *pg.PgDb
	JobId int64
}

func (p *PgOffers) Waiting() []job.Offer {
	//TODO implement me
	panic("implement me")
}

func (pgOffers *PgOffers) Make(ctx context.Context, model job.OfferModel) (job.Offer, error) {
	offerId := int64(0)
	user := user.FetchUser(ctx)

	query := "INSERT INTO job_offer (job_id, price_value, price_currency, description_value, action_created_by_id) VALUES( $1, $2, $3, $4, $5 ) returning id"
	row := pgOffers.Db.Pool.QueryRow(
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
		Db: pgOffers.Db,
		Id: offerId,
	}
	return job.NewSolidOffer(
		&job.OfferModel{
			Id:          offerId,
			Description: model.Description,
			Price:       model.Price,
		},
		&pgOffer,
	), nil
}
