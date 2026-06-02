package pg

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/finance"
)

type PgDacReporting struct {
	db *pg.PgDb
}

func (r *PgDacReporting) isDacReporting() finance.DacReporting {
	return r
}
