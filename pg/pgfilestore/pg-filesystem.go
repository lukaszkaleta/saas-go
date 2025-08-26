package pgfilestoe

import (
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/pg"
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
