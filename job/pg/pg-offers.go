package pgjob

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgOffers struct {
	Db    *pg.PgDb
	JobId int64
}

func (p PgOffers) Waiting() []*job.Offer {
	//TODO implement me
	panic("implement me")
}

func (p PgOffers) Make(person universal.Person, model job.OfferModel) {
	//TODO implement me
	panic("implement me")
}
