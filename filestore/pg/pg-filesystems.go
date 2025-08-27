package pgfilestoe

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
)

type PgFileSystems struct {
	Db *pg.PgDb
}

func NewPgFileSystems(db *pg.PgDb) *PgFileSystems {
	return &PgFileSystems{Db: db}
}
