package pgjob

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgOfferRevisions struct {
	db      *pg.PgDb
	offerId int64
}

func NewPgOfferRevisions(db *pg.PgDb, offerId int64) job.OfferRevisions {
	return &PgOfferRevisions{db: db, offerId: offerId}
}

func (p *PgOfferRevisions) Create(ctx context.Context, model job.OfferRevisionModel) (job.OfferRevision, error) {
	query := "INSERT INTO job_offer_revision (job_offer_id, price_value, price_currency, description_value, action_create_by_id) VALUES( $1, $2, $3, $4, $5 ) returning id"
	var revisionId int64
	err := p.db.Pool.QueryRow(
		ctx,
		query,
		p.offerId,
		model.Price.Value,
		model.Price.Currency,
		model.Description.Value,
		universal.CurrentUserId(ctx),
	).Scan(&revisionId)
	if err != nil {
		return nil, err
	}

	updateQuery := "update job_offer set last_offer_revision_id = $1 where id = $2"
	_, err = p.db.Pool.Exec(ctx, updateQuery, revisionId, p.offerId)
	if err != nil {
		return nil, err
	}

	return p.ById(ctx, revisionId)
}

func (p *PgOfferRevisions) List(ctx context.Context) ([]job.OfferRevision, error) {
	query := "select * from job_offer_revision where job_offer_id = @offerId order by action_created_at desc"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"offerId": p.offerId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOfferRevision(p.db))
}

func (p *PgOfferRevisions) ById(ctx context.Context, id int64) (job.OfferRevision, error) {
	query := "select * from job_offer_revision where id = @id"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOfferRevision(p.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (p *PgOfferRevisions) FromUser(ctx context.Context, id int64) (job.OfferRevision, error) {
	query := "select * from job_offer_revision where job_offer_id = @offerId and action_create_by_id = @userId order by action_created_at desc"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"offerId": p.offerId, "userId": id})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOfferRevision(p.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (p *PgOfferRevisions) NewestFromWorker(ctx context.Context) (job.OfferRevision, error) {
	query := `
		SELECT r.* 
		FROM job_offer_revision r
		JOIN job_offer o ON r.job_offer_id = o.id
		JOIN job j ON o.job_id = j.id
		WHERE r.job_offer_id = @offerId 
		  AND r.action_create_by_id != j.action_created_by_id
		ORDER BY r.action_created_at DESC 
		LIMIT 1`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"offerId": p.offerId})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOfferRevision(p.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (p *PgOfferRevisions) Accepted(ctx context.Context) (job.OfferRevision, error) {
	query := "select * from job_offer_revision where job_offer_id = @offerId and action_accepted_by_id is not null"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"offerId": p.offerId})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOfferRevision(p.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}

func (p *PgOfferRevisions) NewestFromOwner(ctx context.Context) (job.OfferRevision, error) {
	query := `
		SELECT r.* 
		FROM job_offer_revision r
		JOIN job_offer o ON r.job_offer_id = o.id
		JOIN job j ON o.job_id = j.id
		WHERE r.job_offer_id = @offerId 
		  AND r.action_create_by_id = j.action_created_by_id
		ORDER BY r.action_created_at DESC 
		LIMIT 1`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"offerId": p.offerId})
	if err != nil {
		return nil, err
	}
	res, err := pgx.CollectOneRow(rows, MapOfferRevision(p.db))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return res, err
}
