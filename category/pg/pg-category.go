package pgcategory

import (
	"context"

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

func (pgCategory PgCategory) Update(ctx context.Context, model *category.CategoryModel) error {
	return nil
}

func (pgCategory PgCategory) Model(ctx context.Context) *category.CategoryModel {
	return &category.CategoryModel{
		Id: pgCategory.Id,
	}
}

func (pgCategory PgCategory) Localizations() universal.Localizations {
	return &universalPg.PgLocalizations{pgCategory.Db, &pg.RelationEntity{RelationId: pgCategory.Id, ColumnName: "category_id", TableName: "category_localization"}}
}

func (pgCategory PgCategory) Parent(ctx context.Context) (category.Category, error) {
	sql := "select * from category where id = (select parent_category_id from category where id = @id)"
	row, err := pgCategory.Db.Pool.Query(ctx, sql, pgx.NamedArgs{"id": pgCategory.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(row, MapCategoryFunc(pgCategory.Db))
}

func MapCategoryFunc(db *pg.PgDb) pgx.RowToFunc[category.Category] {
	return func(row pgx.CollectableRow) (category.Category, error) {
		return MapCategory(db, row)
	}
}

func MapCategory(db *pg.PgDb, row pgx.CollectableRow) (category.Category, error) {
	model, err := MapCategoryModel(row)
	if err != nil {
		return nil, err
	}
	return category.NewSolidCategory(model, PgCategory{Db: db, Id: model.Id}), nil
}

func MapCategoryModel(row pgx.CollectableRow) (*category.CategoryModel, error) {
	categoryModel := category.EmptyCategoryModel()
	err := row.Scan(
		&categoryModel.Id,
		&categoryModel.ParentId,
		&categoryModel.Name.Value,
		&categoryModel.Name.Slug,
		&categoryModel.Description.Value,
		&categoryModel.Description.ImageUrl)
	if err != nil {
		return categoryModel, err
	}
	return categoryModel, nil
}
