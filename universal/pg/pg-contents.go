package pg

import (
	"context"
	"errors"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

// PgContents implements universal.Contents backed by Postgres content table
// It provides CRUD-ish helpers to fetch and create content rows.
type PgContents struct {
	Db *pg.PgDb
}

func NewPgContents(db *pg.PgDb) universal.Contents {
	return &PgContents{Db: db}
}

// ById loads a content row by id and returns a SolidContent wrapping PgContent
func (p *PgContents) ById(id int64) (universal.Content, error) {
	row := p.Db.Pool.QueryRow(context.Background(), "select id, name_value, name_slug, value from content where id = $1", id)
	model := universal.EmptyContentModel()
	if err := row.Scan(&model.Id, &model.Name.Value, &model.Name.Slug, &model.Value); err != nil {
		return nil, err
	}
	pgContent := &PgContent{Db: p.Db, Id: model.Id}
	return universal.NewSolidContent(model, pgContent), nil
}

// ByName finds the content by name slug or value (tries slug first, then value)
func (p *PgContents) ByName(name string) (universal.Content, error) {
	slug := universal.CreateSlug(name)
	// Try by slug
	row := p.Db.Pool.QueryRow(context.Background(), "select id, name_value, name_slug, value from content where name_slug = $1", slug)
	model := universal.EmptyContentModel()
	if err := row.Scan(&model.Id, &model.Name.Value, &model.Name.Slug, &model.Value); err != nil {
		// try by value as a fallback
		row2 := p.Db.Pool.QueryRow(context.Background(), "select id, name_value, name_slug, value from content where name_value = $1", name)
		if err2 := row2.Scan(&model.Id, &model.Name.Value, &model.Name.Slug, &model.Value); err2 != nil {
			return nil, err2
		}
	}
	pgContent := &PgContent{Db: p.Db, Id: model.Id}
	return universal.NewSolidContent(model, pgContent), nil
}

// Add inserts a new content row and returns Content
func (p *PgContents) Add(model *universal.ContentModel) (universal.Content, error) {
	if model == nil {
		return nil, errors.New("content model is nil")
	}
	// Ensure slug is set
	if model.Name == nil {
		model.Name = universal.EmptyNameModel()
	}
	model.Name.Change(model.Name.Value)
	var id int64
	row, err := p.Db.Pool.Query(context.Background(), "insert into content (name_value, name_slug, value) values ($1, $2, $3) returning id", model.Name.Value, model.Name.Slug, model.Value)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
	}
	model.Id = id
	pgContent := &PgContent{Db: p.Db, Id: id}
	return universal.NewSolidContent(model, pgContent), nil
}
