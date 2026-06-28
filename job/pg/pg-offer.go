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

type PgOffer struct {
	db *pg.PgDb
	Id int64
}

func (pgOffer *PgOffer) ID() int64 {
	return pgOffer.Id
}

func (pgOffer *PgOffer) Create(ctx context.Context, workerId int64) (job.Offer, error) {
	return nil, nil // Not implemented for single offer
}

func (pgOffer *PgOffer) Actions() universal.Actions {
	return pgUniversal.NewPgActions(pgOffer.db, pgOffer.tableEntity())
}

func (pgOffer *PgOffer) Accept(ctx context.Context) error {
	return pgOffer.Actions().WithName(job.Accepted).Execute(ctx)
}

func (pgOffer *PgOffer) Reject(ctx context.Context) error {
	return pgOffer.Actions().WithName(job.Rejected).Execute(ctx)
}

func (pgOffer *PgOffer) Accepted() (bool, error) {
	actionModel := pgOffer.Actions().WithName(job.Accepted).Model(context.Background())
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (pgOffer *PgOffer) Rejected() (bool, error) {
	actionModel := pgOffer.Actions().WithName(job.Rejected).Model(context.Background())
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (pgOffer *PgOffer) Revisions() job.OfferRevisions {
	return NewPgOfferRevisions(pgOffer.db, pgOffer.Id)
}

func (pgOffer *PgOffer) Model(ctx context.Context) (*job.OfferModel, error) {
	query := "select " + OfferColumnString() + " from job_offer where id = @id"
	rows, err := pgOffer.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": pgOffer.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapOfferModel)
}

func OfferColumns() []string {
	return []string{
		"id",
		"job_id",
		"worker_id",
		"status",
		"rating",
		"accepted_offer_revision_id",
		"last_offer_revision_id",
	}
}

func OfferColumnString(alias ...string) string {
	prefix := ""
	if len(alias) > 0 {
		prefix = alias[0] + "."
	}
	columns := OfferColumns()
	aliased := make([]string, len(columns))
	for i, col := range columns {
		aliased[i] = prefix + col
	}
	return strings.Join(aliased, ",")
}

func (pgOffer *PgOffer) tableEntity() pg.TableEntity {
	return pgOffer.db.TableEntity("job_offer", pgOffer.Id)
}

func MapOfferModel(row pgx.CollectableRow) (*job.OfferModel, error) {
	offerModel := job.EmptyOfferModel()
	err := row.Scan(
		&offerModel.Id,
		&offerModel.JobId,
		&offerModel.WorkerId,
		&offerModel.Status,
		&offerModel.Rating,
		&offerModel.AcceptedRevisionId,
		&offerModel.LastRevisionId,
	)
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
