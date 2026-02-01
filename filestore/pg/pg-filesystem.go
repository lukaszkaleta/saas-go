package pgfilestore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgFileSystem struct {
	db    *pg.PgDb
	Id    int64
	Owner pg.RelationEntity
}

func NewPgFileSystem(db *pg.PgDb, id pg.RelationEntity) filestore.FileSystem {
	return &PgFileSystem{
		db:    db,
		Owner: id,
	}
}

func (p PgFileSystem) Model(ctx context.Context) *filestore.FileSystemModel {
	query := "select * from filestore_filesystem where id=&id"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": p.Id})
	if err != nil {
		return nil
	}
	fileSystemModel, err := pgx.CollectOneRow(rows, MapFileSystemModel)
	if err != nil {
		return nil
	}
	return fileSystemModel
}

func (p PgFileSystem) Update(ctx context.Context, newModel *filestore.FileSystemModel) error {
	newModel.Id = p.Id
	query := "update filestore_filesystem set name_value = @nameValue, name_slug = @nameSlug where id = @id"
	cmd, err := p.db.Pool.Exec(ctx, query, FileSystemNamedArgs(newModel))
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (p PgFileSystem) Records() filestore.Records {
	return NewPgRecords(p.db, p.Id)
}

func MapFileSystemModel(row pgx.CollectableRow) (*filestore.FileSystemModel, error) {
	fsm := filestore.EmptyFileSystemModel()
	err := row.Scan(
		&fsm.Id,
		&fsm.Name.Value,
		&fsm.Name.Slug)
	if err != nil {
		return nil, err
	}
	return fsm, nil
}

func MapFileSystem(db *pg.PgDb, owner pg.RelationEntity, row pgx.CollectableRow) (filestore.FileSystem, error) {
	model, err := MapFileSystemModel(row)
	if err != nil {
		return nil, err
	}
	newFileSystem := PgFileSystem{
		db:    db,
		Id:    model.Id,
		Owner: owner,
	}
	return filestore.NewSolidFileSystem(model, newFileSystem), nil
}
