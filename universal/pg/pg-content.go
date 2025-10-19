package pg

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgContent struct {
	Db *pg.PgDb
	Id int64
}

// NewPgContentFromTable constructs PgContent for a specific table entity
func NewPgContentFromTable(db *pg.PgDb, id int64) *PgContent {
	return &PgContent{Db: db, Id: id}
}

func (p *PgContent) Update(model *universal.ContentModel) error {
	// Update name and content value on the owning table
	query := "update content set name_value = $1, name_slug = $2, value = $3 where id = $4"
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Name.Value, model.Name.Slug, model.Value, p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgContent) Model() *universal.ContentModel {
	return &universal.ContentModel{}
}

func (p *PgContent) Localizations() universal.Localizations {
	rel := pg.RelationEntity{TableName: "content_localization", RelationId: p.Id, ColumnName: "content_id"}
	return &PgLocalizations{Db: p.Db, Owner: &rel}
}
