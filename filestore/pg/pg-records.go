package pgfilestoe

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgRecords struct {
	Db *pg.PgDb
}

func (records *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	sql := "insert into filestore_record (name_value, name_slug, description_value, description_image_url) values (@nameValue, @nameSlug, @descriptionValue, @descriptionImageUrl) returning id"
	recordId := int64(0)
	row := records.Db.Pool.QueryRow(ctx, sql, RecordNamedArgs(model))
	err := row.Scan(&recordId)
	if err != nil {
		return nil, err
	}
	newRecord := PgRecord{
		Db: records.Db,
		Id: recordId,
	}
	return filestore.NewSolidRecord(model, newRecord), nil
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
