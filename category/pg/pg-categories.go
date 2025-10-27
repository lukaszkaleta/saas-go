package pgcategory

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/category"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgCategories struct {
	Db *pg.PgDb
}

func NewPgCategories(db *pg.PgDb) *PgCategories {
	return &PgCategories{Db: db}
}

func (pgCategories *PgCategories) AddWithName(nameValue string) (category.Category, error) {
	nameSlug := universal.CreateSlug(nameValue)

	categoryId := int64(0)
	query := "INSERT INTO category(name_value, name_slug) VALUES( $1, $2 ) returning id"
	row := pgCategories.Db.Pool.QueryRow(context.Background(), query, nameValue, nameSlug)
	row.Scan(&categoryId)
	pgCategory := PgCategory{
		Db: pgCategories.Db,
		Id: categoryId,
	}
	return category.NewSolidCategory(
		&category.CategoryModel{
			Id:          categoryId,
			Name:        &universal.NameModel{Value: nameValue, Slug: nameSlug},
			Description: &universal.DescriptionModel{},
		},
		&pgCategory,
	), nil
}

func (pgCategories *PgCategories) AddWithParent(parent category.Category, nameValue string) (category.Category, error) {
	nameSlug := universal.CreateSlug(nameValue)

	categoryId := int64(0)
	query := "INSERT INTO category(parent_category_id, name_value, name_slug) VALUES( $1, $2, $3 ) returning id"
	row := pgCategories.Db.Pool.QueryRow(context.Background(), query, parent.Model().Id, nameValue, nameSlug)
	row.Scan(&categoryId)
	pgCategory := PgCategory{
		Db: pgCategories.Db,
		Id: categoryId,
	}
	return category.NewSolidCategory(
		&category.CategoryModel{
			Id:          categoryId,
			ParentId:    &parent.Model().Id,
			Name:        &universal.NameModel{Value: nameValue, Slug: nameSlug},
			Description: &universal.DescriptionModel{},
		},
		&pgCategory,
	), nil
}

func (pgCategories *PgCategories) AllLocalized(country string, language string) ([]*category.CategoryModel, error) {
	query := "select category_id, (select c.parent_category_id from category c where category_id = id) as parent_category_id, translation_value, translation_slug from category_localization where country = $1 and language = $2"
	rows, err := pgCategories.Db.Pool.Query(context.Background(), query, country, language)
	if err != nil {
		return nil, err
	}

	var categories []*category.CategoryModel
	for rows.Next() {
		categoryModel := new(category.CategoryModel)
		categoryModel.Name = new(universal.NameModel)
		categoryModel.Description = new(universal.DescriptionModel)
		err := rows.Scan(&categoryModel.Id,
			&categoryModel.ParentId,
			&categoryModel.Name.Value,
			&categoryModel.Name.Slug)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categoryModel)
	}
	return categories, nil
}

func (pgCategories *PgCategories) ById(id int64) (category.Category, error) {
	query := "select * from category where id = @id"
	rows, err := pgCategories.Db.Pool.Query(context.Background(), query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	pgCategory := PgCategory{
		Db: pgCategories.Db,
		Id: id,
	}
	categoryModel, err := pgx.CollectOneRow(rows, MapCategoryModel)
	if err != nil {
		return nil, err
	}
	return category.NewSolidCategory(&categoryModel, pgCategory), nil
}

func (pgCategories *PgCategories) ByIds(ids []int64) ([]category.CategoryModel, error) {
	query := "select * from category where id = any($1)"
	rows, err := pgCategories.Db.Pool.Query(context.Background(), query, ids)
	if err != nil {
		return nil, err
	}
	categories, err := pgx.CollectRows(rows, MapCategoryModel)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
