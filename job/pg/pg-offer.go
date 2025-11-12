package pgjob

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgOffer struct {
	Db *pg.PgDb
	Id int64
}

func (o *PgOffer) Accept(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer set action_accepted_at = now(), action_rejected_at = null, action_rejected_by_id = null, action_accepted_by_id = $1 where id = $2"
	_, err := o.Db.Pool.Exec(ctx, query, currentUser.Id, o.Id)
	if err != nil {
		return err
	}
	return nil
}

func (o *PgOffer) Reject(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer set action_rejected_at = now(), action_accepted_at = null, action_accepted_by_id = null ,action_rejected_by_id = $1 where id = $2"
	_, err := o.Db.Pool.Exec(ctx, query, currentUser.Id, o.Id)
	if err != nil {
		return err
	}
	return nil
}

func (o *PgOffer) Accepted() (bool, error) {
	return true, nil
}

func (o *PgOffer) Rejected() (bool, error) {
	return false, nil
}

func (o *PgOffer) Model() *job.OfferModel {
	return &job.OfferModel{}
}

func MapOffer(row pgx.CollectableRow) (*job.OfferModel, error) {
	offerModel := job.EmptyOfferModel()

	nullTimeAccepted := sql.NullTime{}
	nullTimeRejected := sql.NullTime{}

	actionCreated := universal.ActionModel{Name: job.Created}
	actionAccepted := universal.ActionModel{Name: job.Accepted}
	actionRejected := universal.ActionModel{Name: job.Rejected}
	actions := make(map[string]*universal.ActionModel)
	actions[job.Created] = &actionCreated
	actions[job.Accepted] = &actionAccepted
	actions[job.Rejected] = &actionRejected

	err := row.Scan(
		&offerModel.Id,
		&offerModel.JobId,
		&offerModel.Price.Value,
		&offerModel.Price.Currency,
		&offerModel.Description.Value,
		&offerModel.Description.ImageUrl,
		&offerModel.Rating,
		&actionCreated.ById,
		&actionCreated.MadeAt,
		&actionAccepted.ById,
		&nullTimeAccepted,
		&actionRejected.ById,
		&nullTimeRejected,
	)
	actionAccepted.MadeAt = nullTimeAccepted.Time
	actionRejected.MadeAt = nullTimeRejected.Time
	offerModel.Actions = universal.ActionsModel{List: actions}
	if err != nil {
		return nil, err
	}
	return offerModel, nil
}
