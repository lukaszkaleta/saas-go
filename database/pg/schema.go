package pg

import "embed"

type Schema interface {
	Create() error
	CreateTest() error
	Drop() error
	DropTest() error
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

func (schema DefaultSchema) CreateTest() error {
	err := schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/create-test.sql")
	if err != nil {
		return err
	}
	return schema.Create()
}

func (schema DefaultSchema) Drop() error {
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/drop.sql")
}

func (schema DefaultSchema) DropTest() error {
	err := schema.Drop()
	if err != nil {
		return err
	}
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/drop-test.sql")
}
