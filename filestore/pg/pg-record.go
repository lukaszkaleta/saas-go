package pgfilestore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgRecord struct {
	Db *pg.PgDb
	Id int64
}

func (p PgRecord) Model(ctx context.Context) *filestore.RecordModel {
	query := "select * from filestore_record where id=@id"
	rows, err := p.Db.Pool.Query(ctx, query, pgx.NamedArgs{"id": p.Id})
	if err != nil {
		return nil
	}
	recordModel, err := pgx.CollectOneRow(rows, MapRecordModel)
	if err != nil {
		return nil
	}
	return recordModel
}

func (p PgRecord) Update(ctx context.Context, newModel *filestore.RecordModel) error {
	newModel.Id = p.Id
	query := "update filestore_record set name_value = @nameValue, name_slug = @nameSlug where id = @id"
	cmd, err := p.Db.Pool.Exec(ctx, query, RecordNamedArgs(newModel))
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func MapRecordFunc(db *pg.PgDb) pgx.RowToFunc[filestore.Record] {
	return func(row pgx.CollectableRow) (filestore.Record, error) {
		return MapRecord(db, row)
	}
}

func MapRecord(db *pg.PgDb, row pgx.CollectableRow) (filestore.Record, error) {
	model, err := MapRecordModel(row)
	if err != nil {
		return nil, err
	}
	pgMessage := &PgRecord{Db: db, Id: model.Id}
	return filestore.NewSolidRecord(model, pgMessage), nil
}

func MapRecordModel(row pgx.CollectableRow) (*filestore.RecordModel, error) {
	record := filestore.EmptyRecordModel()
	err := row.Scan(
		&record.Id,
		&record.Name.Value,
		&record.Name.Slug,
		&record.Description.Value,
		&record.Description.ImageUrl)
	if err != nil {
		return nil, err
	}
	return record, nil
}
