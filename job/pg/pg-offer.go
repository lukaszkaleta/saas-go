package pgjob

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgOffer struct {
	Db *pg.PgDb
	Id int64
}

func (o *PgOffer) Accept(person universal.Person) error {
	return nil
}

func (o *PgOffer) Reject(person universal.Person) error {
	return nil
}
