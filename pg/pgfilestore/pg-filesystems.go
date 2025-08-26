package pgfilestoe

import (
	"github.com/lukaszkaleta/saas-go/pg/database"
)

type PgFileSystems struct {
	Db *database.PgDb
}

func NewPgFileSystems(db *database.PgDb) *PgFileSystems {
	return &PgFileSystems{Db: db}
}
