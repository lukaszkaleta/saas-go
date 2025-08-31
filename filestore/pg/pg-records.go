package pgfilestore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgRecords struct {
	Db *pg.PgDb
}

func (pg *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	sql := "insert into filestore_record (name_value, name_slug, description_value, description_image_url) values (@nameValue, @nameSlug, @descriptionValue, @descriptionImageUrl) returning id"
	recordId := int64(0)
	row := pg.Db.Pool.QueryRow(ctx, sql, RecordNamedArgs(model))
	err := row.Scan(&recordId)
	if err != nil {
		return nil, err
	}
	newRecord := PgRecord{
		Db: pg.Db,
		Id: recordId,
	}
	model.Id = recordId
	return filestore.NewSolidRecord(model, newRecord), nil
}

func (pg *PgRecords) ById(recordId int64) (filestore.Record, error) {
	sql := "select * from filestore_record where id = @id"
	rows, err := pg.Db.Pool.Query(context.Background(), sql, pgx.NamedArgs{"id": recordId})
	if err != nil {
		return nil, err
	}
	record := PgRecord{
		Db: pg.Db,
		Id: recordId,
	}
	recordModel, err := pgx.CollectOneRow(rows, MapRecord)
	return filestore.NewSolidRecord(recordModel, record), err
}

func RecordNamedArgs(model *filestore.RecordModel) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":                  model.Id,
		"nameValue":           model.Name.Value,
		"nameSlug":            model.Name.Slug,
		"descriptionValue":    model.Description.Value,
		"descriptionImageUrl": model.Description.ImageUrl,
	}
}
