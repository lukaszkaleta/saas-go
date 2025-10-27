package pgcategory

import "github.com/lukaszkaleta/saas-go/database/pg"

type PgRelationCategories struct {
	Db         *pg.PgDb
	Categories *PgCategories
	relation   pg.RelationEntity
}

func NewPgRelationCategories(pgCategories *PgCategories, relation pg.RelationEntity) PgRelationCategories {
	return PgRelationCategories{
		Db:         pgCategories.Db,
		relation:   relation,
		Categories: pgCategories,
	}
}
