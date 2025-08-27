package pg

import "embed"

type Schema interface {
	Create() error
	Drop() error
}

type DefaultSchema struct {
	db    *PgDb
	ddlFs embed.FS
}

func NewDefaultSchema(db *PgDb, ddlFs embed.FS) Schema {
	return &DefaultSchema{db: db, ddlFs: ddlFs}
}

func (schema DefaultSchema) Create() error {
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/create.sql")
}

func (schema DefaultSchema) Drop() error {
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/drop.sql")
}
