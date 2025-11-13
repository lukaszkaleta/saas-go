package pgjob

import (
	"embed"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

//go:embed ddl
var ddlFs embed.FS

func NewJobSchema(db *pg.PgDb) pg.Schema {
	return pg.NewDefaultSchema(db, ddlFs, "job")
}
