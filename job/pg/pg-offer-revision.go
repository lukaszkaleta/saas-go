package pgjob

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
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

func (p *PgOfferRevision) Accept(ctx context.Context) error {
	err := p.Actions().WithName(job.Accepted).Execute(ctx)
	if err != nil {
		return err
	}

	model, err := p.Model(ctx)
	if err != nil {
		return err
	}

	query := "update job_offer set accepted_offer_revision_id = @revisionId, last_offer_revision_id = @revisionId, status = @status where id = @offerId"
	_, err = p.db.Pool.Exec(ctx, query, pgx.NamedArgs{
		"revisionId": p.Id,
		"status":     job.Accepted,
		"offerId":    model.OfferId,
	})

	return err
}

func (p *PgOfferRevision) Reject(ctx context.Context) error {
	err := p.Actions().WithName(job.Rejected).Execute(ctx)
	if err != nil {
		return err
	}

	model, err := p.Model(ctx)
	if err != nil {
		return err
	}

	query := "update job_offer set status = @status where id = @offerId"
	_, err = p.db.Pool.Exec(ctx, query, pgx.NamedArgs{
		"status":  job.Rejected,
		"offerId": model.OfferId,
	})

	return err
}

func (pgRevision *PgOfferRevision) Accepted() (bool, error) {
	actionModel := pgRevision.Actions().WithName(job.Accepted).Model(context.Background())
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (pgRevision *PgOfferRevision) Rejected() (bool, error) {
	actionModel := pgRevision.Actions().WithName(job.Rejected).Model(context.Background())
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (pgRevision *PgOfferRevision) Model(ctx context.Context) (*job.OfferRevisionModel, error) {
	query := "select " + OfferRevisionColumnString() + " from job_offer_revision where id = @id"
	rows, err := pgRevision.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": pgRevision.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOfferRevisionModel)
}

func (pgRevision *PgOfferRevision) tableEntity() pg.TableEntity {
	return pgRevision.db.TableEntity("job_offer_revision", pgRevision.Id)
}

func OfferRevisionColumns() []string {
	return []string{
		"id",
		"job_offer_id",
		"price_value",
		"price_currency",
		"description_value",
		"description_image_url",
		"action_created_by_id",
		"action_created_at",
		"action_accepted_by_id",
		"action_accepted_at",
		"action_rejected_by_id",
		"action_rejected_at",
	}
}

func OfferRevisionColumnString(alias ...string) string {
	prefix := ""
	if len(alias) > 0 {
		prefix = alias[0] + "."
	}
	columns := OfferRevisionColumns()
	aliased := make([]string, len(columns))
	for i, col := range columns {
		aliased[i] = prefix + col
	}
	return strings.Join(aliased, ",")
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
