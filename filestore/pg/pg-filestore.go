package pgfilestoe

import (
	"embed"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

//go:embed ddl
var ddlFs embed.FS

type FilestoreSchema struct {
	Db *pg.PgDb
}

func (filestoreSchema FilestoreSchema) Create() error {
	return filestoreSchema.Db.ExecuteFileFromFs(ddlFs, "ddl/create.sql")
}

func (filestoreSchema FilestoreSchema) Drop() error {
	return filestoreSchema.Db.ExecuteFileFromFs(ddlFs, "ddl/drop.sql")
}
