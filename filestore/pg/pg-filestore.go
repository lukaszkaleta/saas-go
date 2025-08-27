package pgfilestoe

import "github.com/lukaszkaleta/saas-go/database/pg"

type FilestoreSchema struct {
	Db *pg.PgDb
}

func (filestoreSchema FilestoreSchema) Create() error {
	return filestoreSchema.Db.ExecuteSql("ddl.sql")
}

func (filestoreSchema FilestoreSchema) Drop() error {
	return filestoreSchema.Db.ExecuteSql("drop.sql")
}
