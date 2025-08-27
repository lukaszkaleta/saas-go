package pgfilestoe

import (
	"embed"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

var folder embed.FS

type FilestoreSchema struct {
	Db *pg.PgDb
}

func (filestoreSchema FilestoreSchema) Create() error {
	return filestoreSchema.Db.ExecuteFileFromFs(folder, "ddl.sql")
}

func (filestoreSchema FilestoreSchema) Drop() error {
	return filestoreSchema.Db.ExecuteFileFromFs(folder, "drop.sql")
}
