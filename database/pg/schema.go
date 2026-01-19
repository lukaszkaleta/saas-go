package pg

import (
	"embed"
)

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
	err := NewSqlFile(schema.ddlFs, schema.db, "ddl/create.sql").Execute()
	if err != nil {
		return err
	}
	err = NewSqlFile(schema.ddlFs, schema.db, "ddl/create-func.sql", WithSkipNotExistingFile, WithQuerySeparator(";;")).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (schema *DefaultSchema) CreateTest() error {
	err := NewSqlFile(schema.ddlFs, schema.db, "ddl/create-test.sql", WithSkipNotExistingFile).Execute()
	if err != nil {
		return err
	}
	return schema.Create()
}

func (schema *DefaultSchema) Drop() error {
	err := NewSqlFile(schema.ddlFs, schema.db, "ddl/drop.sql").Execute()
	if err != nil {
		return err
	}
	return NewSqlFile(schema.ddlFs, schema.db, "ddl/drop-func.sql", WithSkipNotExistingFile).Execute()
}

func (schema *DefaultSchema) DropTest() error {
	err := schema.Drop()
	if err != nil {
		return err
	}
	return NewSqlFile(schema.ddlFs, schema.db, "ddl/drop-test.sql", WithSkipNotExistingFile).Execute()
}
