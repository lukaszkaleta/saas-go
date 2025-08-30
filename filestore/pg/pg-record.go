package pgfilestoe

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	pguniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgRecord struct {
	Db *pg.PgDb
	Id int64
}

func (p PgRecord) Model() *filestore.RecordModel {
	query := "select * from record where id=@id"
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
	query := "update filestore_record set name_value = @nameValue, nameSlug = @nameSlug where id = @id"
	_, err := p.Db.Pool.Exec(context.Background(), query, RecordNamedArgs(newModel))
	if err != nil {
		return err
	}
	return nil
}

func MapRecord(row pgx.CollectableRow) (*filestore.RecordModel, error) {
	record := filestore.EmptyRecordModel()
	pguniversal.UseMapName(record.Name)(row)
	pguniversal.UseMapDescription(record.Description)(row)
	return record, nil
}
