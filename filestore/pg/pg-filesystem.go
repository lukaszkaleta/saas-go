package pgfilestoe

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgFileSystem struct {
	Db    *pg.PgDb
	Owner pg.RelationEntity
}

func (p *PgFileSystem) Update(model *filestore.FileSystemModel) error {
	return nil
}

func (p *PgFileSystem) Model() *filestore.FileSystemModel {
	return &filestore.FileSystemModel{}
}

func (p *PgFileSystem) Records() *filestore.Records {
	return nil
}
