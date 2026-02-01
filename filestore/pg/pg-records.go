package pgfilestore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRecords struct {
	db           *pg.PgDb
	filesystemId int64
}

func (pg *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	sql := "insert into filestore_record (name_value, name_slug, description_value, description_image_url) values (@nameValue, @nameSlug, @descriptionValue, @descriptionImageUrl) returning id"
	recordId := int64(0)
	row := pg.db.Pool.QueryRow(ctx, sql, RecordNamedArgs(model))
	err := row.Scan(&recordId)
	if err != nil {
		return nil, err
	}
	newRecord := PgRecord{
		Db: pg.db,
		Id: recordId,
	}
	model.Id = recordId
	return filestore.NewSolidRecord(model, newRecord), nil
}

func (pg *PgRecords) AddFromUrl(ctx context.Context, url string) (filestore.Record, error) {
	model := filestore.EmptyRecordModel()
	model.Url = url
	return pg.Add(ctx, model)
}

func (pg *PgRecords) AddAll(ctx context.Context, models []*filestore.RecordModel) ([]filestore.Record, error) {
	created := make([]filestore.Record, len(models))
	for i, record := range models {
		added, err := pg.Add(ctx, record)
		if err != nil {
			return nil, err
		}
		created[i] = added
	}
	return created, nil
}

func (pg *PgRecords) AddFromUrls(ctx context.Context, urls []string) ([]filestore.Record, error) {
	created := make([]filestore.Record, len(urls))
	for i, url := range urls {
		added, err := pg.AddFromUrl(ctx, url)
		if err != nil {
			return nil, err
		}
		created[i] = added
	}
	return created, nil
}

func (pg *PgRecords) AddFromName(ctx context.Context, name string) (filestore.Record, error) {
	model := filestore.EmptyRecordModel()
	model.Name = universal.SluggedName(name)
	return pg.Add(ctx, model)
}

func (pg *PgRecords) ById(ctx context.Context, recordId int64) (filestore.Record, error) {
	sql := "select * from filestore_record where id = @id"
	rows, err := pg.db.Pool.Query(ctx, sql, pgx.NamedArgs{"id": recordId})
	if err != nil {
		return nil, err
	}
	record := PgRecord{
		Db: pg.db,
		Id: recordId,
	}
	recordModel, err := pgx.CollectOneRow(rows, MapRecordModel)
	return filestore.NewSolidRecord(recordModel, record), err
}

func (pg *PgRecords) Urls(ctx context.Context) ([]string, error) {
	//sql := "select description_image_url from filestore_record where description_image_url = @url"
	return make([]string, 0), nil
}

func (pg *PgRecords) FindByUrl(ctx context.Context, url string) (filestore.Record, error) {
	sql := "select * from filestore_record where description_image_url = @url"
	rows, err := pg.db.Pool.Query(ctx, sql, pgx.NamedArgs{"url": url})
	if err != nil {
		return nil, err
	}
	return MapRecord(pg.db, rows)
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
