package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgOffer struct {
	db *pg.PgDb
	Id int64
}

func (pgOffer *PgOffer) Accept(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer set action_accepted_at = now(), action_rejected_at = null, action_rejected_by_id = null, action_accepted_by_id = $1 where id = $2"
	_, err := pgOffer.db.Pool.Exec(ctx, query, currentUser.Id, pgOffer.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pgOffer *PgOffer) Reject(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer set action_rejected_at = now(), action_accepted_at = null, action_accepted_by_id = null ,action_rejected_by_id = $1 where id = $2"
	_, err := pgOffer.db.Pool.Exec(ctx, query, currentUser.Id, pgOffer.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pgOffer *PgOffer) Accepted() (bool, error) {
	return true, nil
}

func (pgOffer *PgOffer) Rejected() (bool, error) {
	return false, nil
}

func (pgOffer *PgOffer) Model(ctx context.Context) (*job.OfferModel, error) {
	query := "select * from job_offer where job_id = @jobId"
	rows, err := pgOffer.db.Pool.Query(ctx, query, pgx.NamedArgs{"jobId": pgOffer.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOfferModel)
}

func MapOfferModel(row pgx.CollectableRow) (*job.OfferModel, error) {
	offerModel := job.EmptyOfferModel()
	actionCreatedModel := universal.EmptyCreatedActionModel()
	actionAcceptedModel := universal.EmptyActionModel(job.Accepted)
	actionRejectedModel := universal.EmptyActionModel(job.Rejected)
	err := row.Scan(
		&offerModel.Id,
		&offerModel.JobId,
		&offerModel.Price.Value,
		&offerModel.Price.Currency,
		&offerModel.Description.Value,
		&offerModel.Description.ImageUrl,
		&offerModel.Rating,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
		&actionAcceptedModel.ById,
		&actionAcceptedModel.MadeAt,
		&actionRejectedModel.ById,
		&actionRejectedModel.MadeAt,
	)
	offerModel.Actions.List[actionCreatedModel.Name] = actionCreatedModel
	offerModel.Actions.List[actionAcceptedModel.Name] = actionCreatedModel
	offerModel.Actions.List[actionRejectedModel.Name] = actionCreatedModel
	if err != nil {
		return nil, err
	}
	return offerModel, nil

}

func MapOffer(db *pg.PgDb) pgx.RowToFunc[job.Offer] {
	return func(row pgx.CollectableRow) (job.Offer, error) {
		model, err := MapOfferModel(row)
		if err != nil {
			return nil, err
		}
		pgOffer := &PgOffer{db: db, Id: model.Id}
		solidOffer := job.NewSolidOffer(model, pgOffer)
		return solidOffer, nil
	}
}
