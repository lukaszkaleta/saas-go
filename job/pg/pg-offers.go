package pgjob

import (
	"context"

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

func (pgOffers *PgOffers) ById(ctx context.Context, id int64) (job.Offer, error) {
	query := "select * from job_offer where job_id = @jobId and id = @id"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffers.JobId, "id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOffer(pgOffers.db))
}

func (pgOffers *PgOffers) Waiting(ctx context.Context) ([]job.Offer, error) {
	query := "select * from job_offer where job_id = $1 and action_accepted_at is null and action_rejected_at is null"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgOffers.JobId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOffer(pgOffers.db))
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
	actionsList[job.Created] = universal.NowActionModelForUser(job.Created, &user.Id)
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

func (pgOffers *PgOffers) FromUser(ctx context.Context, user universal.Idable) (job.Offer, error) {
	query := "select * from job_offer where job_id = @jobId and action_created_by_id = @userId order by action_created_at_desc limit 1"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffers.JobId, "userId": user.ID()})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOffer(pgOffers.db))
}
