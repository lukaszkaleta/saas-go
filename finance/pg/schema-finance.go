package pg

import (
	"embed"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

//go:embed ddl
var ddlFs embed.FS

func NewFinanceSchema(db *pg.PgDb) pg.Schema {
	return pg.NewDefaultSchema(db, ddlFs, "finance")
}
