package pgcategory

import (
	"embed"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

//go:embed ddl
var ddlFs embed.FS

func NewCategorySchema(db *pg.PgDb) pg.Schema {
	return pg.NewDefaultSchema(db, ddlFs)
}
