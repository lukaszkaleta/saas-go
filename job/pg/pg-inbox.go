package pgjob

import (
	"context"

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
	//TODO implement me
	panic("implement me")
}

func (p PgOfferInbox) CountUnread(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func NewPgOfferInbox(db *pg.PgDb) universal.Inbox[job.Offer] {
	return PgOfferInbox{db: db}
}

type PgTaskInbox struct {
	db *pg.PgDb
}
