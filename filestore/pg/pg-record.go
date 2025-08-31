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

func (p PgRecord) Model() *filestore.RecordModel {
	query := "select * from filestore_record where id=@id"
	rows, err := p.Db.Pool.Query(context.Background(), query, pgx.NamedArgs{"id": p.Id})
	if err != nil {
		return nil
	}
	recordModel, err := pgx.CollectOneRow(rows, MapRecord)
	if err != nil {
		return nil
	}
	return recordModel
}

func (p PgRecord) Update(newModel *filestore.RecordModel) error {
	newModel.Id = p.Id
	query := "update filestore_record set name_value = @nameValue, name_slug = @nameSlug where id = @id"
	cmd, err := p.Db.Pool.Exec(context.Background(), query, RecordNamedArgs(newModel))
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func MapRecord(row pgx.CollectableRow) (*filestore.RecordModel, error) {
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
