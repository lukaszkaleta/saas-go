package pgcategory

import (
	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/category"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	universalPg "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgCategory struct {
	Db *pg.PgDb
	Id int64
}

func (pgCategory PgCategory) Update(model *category.CategoryModel) error {
	return nil
}

func (pgCategory PgCategory) Model() *category.CategoryModel {
	return &category.CategoryModel{
		Id: pgCategory.Id,
	}
}

func (pgCategory PgCategory) Localizations() universal.Localizations {
	return &universalPg.PgLocalizations{pgCategory.Db, &pg.RelationEntity{RelationId: pgCategory.Id, ColumnName: "category_id", TableName: "category_localization"}}
}

func MapCategory(db *pg.PgDb) pgx.RowToFunc[category.Category] {
	return func(row pgx.CollectableRow) (category.Category, error) {
		model, err := MapCategoryModel(row)
		if err != nil {
			return nil, err
		}
		return category.NewSolidCategory(&model, PgCategory{Db: db, Id: model.Id}), nil
	}
}

func MapCategoryModel(row pgx.CollectableRow) (category.CategoryModel, error) {
	categoryModel := category.EmptyCategoryModel()
	err := row.Scan(
		&categoryModel.Id,
		&categoryModel.ParentId,
		&categoryModel.Name.Value,
		&categoryModel.Name.Slug,
		&categoryModel.Description.Value,
		&categoryModel.Description.ImageUrl)
	if err != nil {
		return *categoryModel, err
	}
	return *categoryModel, nil
}
