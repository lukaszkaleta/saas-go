package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgOfferRevision struct {
	db *pg.PgDb
	Id int64
}

func (pgRevision *PgOfferRevision) ID() int64 {
	return pgRevision.Id
}

func (pgRevision *PgOfferRevision) Actions() universal.Actions {
	return pgUniversal.NewPgActions(pgRevision.db, pgRevision.tableEntity())
}

func (pgRevision *PgOfferRevision) Accept(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer_revision set action_accepted_at = now(), action_accepted_by_id = @userId where id = @id"
	_, err := pgRevision.db.Pool.Exec(ctx, query, pgx.NamedArgs{"id": pgRevision.Id, "userId": currentUser.Id})
	if err != nil {
		return err
	}
	return nil
}

func (pgRevision *PgOfferRevision) Reject(ctx context.Context) error {
	currentUser := user.CurrentUser(ctx)
	query := "update job_offer_revision set action_rejected_at = now(), action_rejected_by_id = @userId where id = @id"
	_, err := pgRevision.db.Pool.Exec(ctx, query, pgx.NamedArgs{"id": pgRevision.Id, "userId": currentUser.Id})
	if err != nil {
		return err
	}
	return nil
}

func (pgRevision *PgOfferRevision) Accepted() (bool, error) {
	query := "select action_accepted_by_id is not null from job_offer_revision where id = @id"
	var accepted bool
	err := pgRevision.db.Pool.QueryRow(context.Background(), query, pgx.NamedArgs{"id": pgRevision.Id}).Scan(&accepted)
	if err != nil {
		return false, err
	}
	return accepted, nil
}

func (pgRevision *PgOfferRevision) Rejected() (bool, error) {
	query := "select action_rejected_by_id is not null from job_offer_revision where id = @id"
	var rejected bool
	err := pgRevision.db.Pool.QueryRow(context.Background(), query, pgx.NamedArgs{"id": pgRevision.Id}).Scan(&rejected)
	if err != nil {
		return false, err
	}
	return rejected, nil
}

func (pgRevision *PgOfferRevision) Model(ctx context.Context) (*job.OfferRevisionModel, error) {
	query := "select * from job_offer_revision where id = @id"
	rows, err := pgRevision.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": pgRevision.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOfferRevisionModel)
}

func (pgRevision *PgOfferRevision) tableEntity() pg.TableEntity {
	return pgRevision.db.TableEntity("job_offer_revision", pgRevision.Id)
}

func MapOfferRevisionModel(row pgx.CollectableRow) (*job.OfferRevisionModel, error) {
	model := job.EmptyOfferRevisionModel()
	actionCreatedModel := universal.EmptyCreatedActionModel()
	actionAcceptedModel := universal.EmptyActionModel(job.Accepted)
	actionRejectedModel := universal.EmptyActionModel(job.Rejected)

	err := row.Scan(
		&model.Id,
		&model.OfferId,
		&model.Price.Value,
		&model.Price.Currency,
		&model.Description.Value,
		&model.Description.ImageUrl,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
		&actionAcceptedModel.ById,
		&actionAcceptedModel.MadeAt,
		&actionRejectedModel.ById,
		&actionRejectedModel.MadeAt,
	)
	if err != nil {
		return nil, err
	}
	model.Actions.List[actionCreatedModel.Name] = actionCreatedModel
	model.Actions.List[actionAcceptedModel.Name] = actionAcceptedModel
	model.Actions.List[actionRejectedModel.Name] = actionRejectedModel

	return model, nil
}

func MapOfferRevision(db *pg.PgDb) pgx.RowToFunc[job.OfferRevision] {
	return func(row pgx.CollectableRow) (job.OfferRevision, error) {
		model, err := MapOfferRevisionModel(row)
		if err != nil {
			return nil, err
		}
		pgRevision := &PgOfferRevision{db: db, Id: model.Id}
		solidRevision := job.NewSolidOfferRevision(model, pgRevision)
		return solidRevision, nil
	}
}
