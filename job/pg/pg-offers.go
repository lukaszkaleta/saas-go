package pgjob

import (
	"context"
	"errors"

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
	query := "select " + OfferColumnString() + " from job_offer where job_id = @jobId and id = @id"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffers.JobId, "id": id})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOffer(pgOffers.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (pgOffers *PgOffers) Waiting(ctx context.Context) ([]job.Offer, error) {
	query := "select " + OfferColumnString() + " from job_offer where job_id = $1 and status = 'waiting'"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgOffers.JobId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOffer(pgOffers.db))
}

func (pgOffers *PgOffers) Accepted(ctx context.Context) (job.Offer, error) {
	query := "select " + OfferColumnString() + " from job_offer where job_id = $1 and accepted_offer_revision_id is not null"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgOffers.JobId)
	if err != nil {
		return nil, err
	}
	collectRows, err := pgx.CollectRows(rows, MapOffer(pgOffers.db))
	if err != nil {
		return nil, err
	}
	if len(collectRows) == 0 {
		return nil, nil
	}
	return collectRows[0], nil
}

func (pgOffers *PgOffers) Make(ctx context.Context, workerId int64, model *job.OfferRevisionModel) (job.OfferRevision, error) {
	offerFromUser, err := pgOffers.FromUser(ctx, universal.JustId{Id: workerId})
	if err != nil {
		return nil, err
	}
	if offerFromUser == nil {
		offerFromUser, err = pgOffers.Create(ctx, workerId)
		if err != nil {
			return nil, err
		}
	}

	return offerFromUser.Revisions().Create(ctx, *model)
}

func (pgOffers *PgOffers) FromUser(ctx context.Context, user universal.Idable) (job.Offer, error) {
	query := "select " + OfferColumnString() + " from job_offer where job_id = @jobId and worker_id = @userId limit 1"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffers.JobId, "userId": user.ID()})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOffer(pgOffers.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (pgOffers *PgOffers) Delete(ctx context.Context) error {
	query := "DELETE FROM job_offer WHERE job_id = $1"
	_, err := pgOffers.db.Pool.Exec(ctx, query, pgOffers.JobId)
	return err
}

func (pgOffers *PgOffers) waitingOfferFromWorker(ctx context.Context, user *user.UserModel) (job.Offer, error) {
	query := "select " + OfferColumnString() + " from job_offer where job_id = @jobId and action_accepted_at is null and worker_id = @userId limit 1"
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffers.JobId, "userId": user.ID()})
	if err != nil {
		return nil, err
	}
	offers, err := pgx.CollectRows(rows, MapOffer(pgOffers.db))
	if err != nil {
		return nil, err
	}
	if len(offers) == 0 {
		return nil, nil
	}
	return offers[0], nil
}

func (pgOffers *PgOffers) Create(ctx context.Context, workerId int64) (job.Offer, error) {
	query := "INSERT INTO job_offer (job_id, worker_id) VALUES ($1, $2) RETURNING " + OfferColumnString()
	rows, err := pgOffers.db.Pool.Query(ctx, query, pgOffers.JobId, workerId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOffer(pgOffers.db))
}
