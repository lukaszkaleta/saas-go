package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/messages"
	pgMessages "github.com/lukaszkaleta/saas-go/messages/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgJobInbox struct {
	db *pg.PgDb
}

func NewPgJobInbox(db *pg.PgDb) *PgJobInbox {
	return &PgJobInbox{db: db}
}

func (p *PgJobInbox) Messages() universal.Inbox[messages.Message] {
	return pgMessages.NewPgQuestionInbox(p.db, pg.RelationEntity{})
}

func (p *PgJobInbox) Offers() universal.Inbox[job.Offer] {
	return NewPgOfferInbox(p.db)
}

type PgOfferInbox struct {
	db *pg.PgDb
}

func (p PgOfferInbox) Last(ctx context.Context) ([]job.Offer, error) {
	query := `
		SELECT jo.* 
		FROM job_offer jo
		JOIN job j ON jo.job_id = j.id
		WHERE j.action_created_by_id = @userId 
		  AND jo.action_accepted_at IS NULL 
		  AND jo.action_rejected_at IS NULL
		ORDER BY jo.action_created_at DESC
	`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": universal.CurrentUserId(ctx)})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapOffer(p.db))
}

func (p PgOfferInbox) CountUnread(ctx context.Context) (int, error) {
	query := `
		SELECT count(*) 
		FROM job_offer jo
		JOIN job j ON jo.job_id = j.id
		WHERE j.action_created_by_id = @userId 
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

func NewPgOfferInbox(db *pg.PgDb) universal.Inbox[job.Offer] {
	return PgOfferInbox{db: db}
}

type PgTaskInbox struct {
	db *pg.PgDb
}
