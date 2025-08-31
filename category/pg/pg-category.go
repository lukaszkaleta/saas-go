package pgcategory

import (
	"github.com/lukaszkaleta/saas-go/category"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	universalPg "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgCategory struct {
	Db *pg.PgDb
	Id int64
}

func (pgCategory *PgCategory) Update(model *category.CategoryModel) error {
	return nil
}

func (pgCategory *PgCategory) Model() *category.CategoryModel {
	return &category.CategoryModel{
		Id: pgCategory.Id,
	}
}

func (pgCategory *PgCategory) Localizations() universal.Localizations {
	return &universalPg.PgLocalizations{pgCategory.Db, &pg.RelationEntity{RelationId: pgCategory.Id, ColumnName: "category_id", TableName: "category_localization"}}
}
