package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgJobOutbox struct {
	db *pg.PgDb
}

func NewPgJobOutbox(db *pg.PgDb) *PgJobOutbox {
	return &PgJobOutbox{db: db}
}

func (p *PgJobOutbox) Offers() universal.Outbox[job.Offer] {
	return NewPgOfferOutbox(p.db)
}

type PgOfferOutbox struct {
	db *pg.PgDb
}

func (p PgOfferOutbox) Last(ctx context.Context) ([]job.Offer, error) {
	query := `
		SELECT jo.* 
		FROM job_offer jo
		JOIN job j ON jo.job_id = j.id
		WHERE jo.action_created_by_id = @userId
		AND jo.action_accepted_at IS NULL 
		AND jo.action_rejected_at IS NULL
		ORDER BY jo.action_created_at DESC
	`
	currentUserId := universal.CurrentUserId(ctx)
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": currentUserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOffer(p.db))
}

func (p PgOfferOutbox) Old(ctx context.Context) ([]job.Offer, error) {
	query := `
		SELECT jo.* 
		FROM job_offer jo
		JOIN job j ON jo.job_id = j.id
		WHERE jo.action_created_by_id = @userId 
		AND (jo.action_accepted_at IS NOT NULL OR jo.action_rejected_at IS NOT NULL)
		ORDER BY jo.action_created_at DESC
	`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": universal.CurrentUserId(ctx)})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOffer(p.db))
}

func (p PgOfferOutbox) CountUnread(ctx context.Context) (int, error) {
	query := `
		SELECT count(*) 
		FROM job_offer jo
		JOIN job j ON jo.job_id = j.id
		WHERE jo.action_created_by_id = @userId 
		  AND jo.action_accepted_at IS NULL 
		  AND jo.action_rejected_at IS NULL
	`
	var count int
	err := p.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"userId": universal.CurrentUserId(ctx)}).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewPgOfferOutbox(db *pg.PgDb) universal.Outbox[job.Offer] {
	return PgOfferOutbox{db: db}
}
