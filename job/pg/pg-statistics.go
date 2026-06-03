package pgjob

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgStatistics struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (s *PgStatistics) Clicks() universal.Counter {
	return &pgUniversal.PgCounter{
		Db:          s.Db,
		TableEntity: s.TableEntity,
		ColumnName:  "clicks",
	}
}
