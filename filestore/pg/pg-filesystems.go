package pgfilestoe

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgFileSystems struct {
	Db *pg.PgDb
}

func NewPgFileSystems(db *pg.PgDb) *PgFileSystems {
	return &PgFileSystems{Db: db}
}

func (pgFileSystem PgFileSystems) Add(ctx context.Context, name string, ownerId int64) (filestore.FileSystem, error) {
	sql := "insert into filestore_filesystem (name_value, name_slug) values (@nameValue, @nameSlug) returning id"
	newId := int64(0)
	model := &filestore.FileSystemModel{Name: universal.SluggedName(name)}
	row := pgFileSystem.Db.Pool.QueryRow(ctx, sql, FileSystemNamedArgs(model))
	err := row.Scan(&newId)
	if err != nil {
		return nil, err
	}
	newFileSystem := PgFileSystem{
		Db:    pgFileSystem.Db,
		Id:    newId,
		Owner: FileSystemRelationEntity(name, ownerId),
	}
	model.Id = newId
	return filestore.NewSolidFileSystem(model, newFileSystem), nil
}

func FileSystemOwnerColumnName(name string) string {
	return name + "_id"
}

func FileSystemOwnerTableName(name string) string {
	return name + "_filesystem"
}

func FileSystemRelationEntity(name string, ownerId int64) pg.RelationEntity {
	return pg.RelationEntity{RelationId: ownerId, TableName: FileSystemOwnerTableName(name), ColumnName: FileSystemOwnerColumnName(name)}
}

func FileSystemNamedArgs(model *filestore.FileSystemModel) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":        model.Id,
		"nameValue": model.Name.Value,
		"nameSlug":  model.Name.Slug,
	}
}
