package postgres

import "github.com/lukaszkaleta/saas-go/pg"

type PgFileSystems struct {
	Db *pg.PgDb
}

func NewPgFileSystems(db *pg.PgDb) *PgFileSystems {
	return &PgFileSystems{Db: db}
}
