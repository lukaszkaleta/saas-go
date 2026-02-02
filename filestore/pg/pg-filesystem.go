package pgfilestore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgFileSystem struct {
	db    *pg.PgDb
	Id    int64
	Owner pg.RelationEntity
}

func NewPgFileSystem(db *pg.PgDb, owner pg.RelationEntity) filestore.FileSystem {
	return &PgFileSystem{
		db:    db,
		Owner: owner,
	}
}

func (p *PgFileSystem) ID() int64 {
	return p.Id
}

func (p *PgFileSystem) Model(ctx context.Context) *filestore.FileSystemModel {
	query := "select * from filestore_filesystem where id=&id"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": p.ID()})
	if err != nil {
		return nil
	}
	fileSystemModel, err := pgx.CollectOneRow(rows, MapFileSystemModel)
	if err != nil {
		return nil
	}
	return fileSystemModel
}

func (p *PgFileSystem) Update(ctx context.Context, newModel *filestore.FileSystemModel) error {
	newModel.Id = p.ID()
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

func (p *PgFileSystem) Records() filestore.Records {
	return NewPgRecords(p.db, p)
}

func (p *PgFileSystem) Init(ctx context.Context) (int64, error) {

	sql := fmt.Sprintf("select filesystem_id from %s where %s = @relationId", p.Owner.TableName, p.Owner.ColumnName)
	rows, err := p.db.Pool.Query(ctx, sql, pgx.NamedArgs{"relationId": p.Owner.RelationId})
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	filesystemId := int64(0)
	if rows.Next() {
		err := rows.Scan(&filesystemId)
		if err != nil {
			return 0, err
		}
		if filesystemId > 0 {
			p.Id = filesystemId
			return filesystemId, nil
		}
	}

	sql = "insert into filestore_filesystem (name_value, name_slug) values (@name, @slug) returning id"
	row := p.db.Pool.QueryRow(ctx, sql, pgx.NamedArgs{"name": p.Owner.TableName, "slug": p.Owner.TableName})
	err = row.Scan(&filesystemId)
	if err != nil {
		return 0, err
	}
	sql = fmt.Sprintf("insert into %s (filesystem_id, %s) values (@filesystemId, @relationId)", p.Owner.TableName, p.Owner.ColumnName)
	_, err = p.db.Pool.Exec(ctx, sql, pgx.NamedArgs{"filesystemId": filesystemId, "relationId": p.Owner.RelationId})
	if err != nil {
		return 0, err
	}

	p.Id = filesystemId
	return filesystemId, nil
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
	newFileSystem := &PgFileSystem{
		db:    db,
		Id:    model.Id,
		Owner: owner,
	}
	return filestore.NewSolidFileSystem(model, newFileSystem), nil
}
