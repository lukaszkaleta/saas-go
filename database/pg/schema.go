package pg

import "embed"

type Schema interface {
	Name() string
	Create() error
	CreateTest() error
	Drop() error
	DropTest() error
}

type DefaultSchema struct {
	db    *PgDb
	ddlFs embed.FS
	name  string
}

func NewDefaultSchema(db *PgDb, ddlFs embed.FS, name string) Schema {
	return &DefaultSchema{db: db, ddlFs: ddlFs, name: name}
}

func (s *DefaultSchema) Name() string {
	return s.name
}

func (schema *DefaultSchema) Create() error {
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/create.sql")
}

func (schema *DefaultSchema) CreateTest() error {
	err := schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/create-test.sql")
	if err != nil {
		return err
	}
	return schema.Create()
}

func (schema *DefaultSchema) Drop() error {
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/drop.sql")
}

func (schema *DefaultSchema) DropTest() error {
	err := schema.Drop()
	if err != nil {
		return err
	}
	return schema.db.ExecuteFileFromFs(schema.ddlFs, "ddl/drop-test.sql")
}
