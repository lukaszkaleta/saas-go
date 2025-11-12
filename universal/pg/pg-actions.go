package pg

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func NewPgActions(db *pg.PgDb, tableEntity pg.TableEntity) universal.Actions {
	return &PgActions{db: db, tableEntity: tableEntity}
}

type PgActions struct {
	db          *pg.PgDb
	tableEntity pg.TableEntity
}

func (p PgActions) WithName(name string) universal.Action {
	//TODO implement me
	panic("implement me")
}

func (p PgActions) List() map[string]*universal.Action {
	//TODO implement me
	panic("implement me")
}

func (p PgActions) Model() *universal.ActionsModel {
	//TODO implement me
	panic("implement me")
}
