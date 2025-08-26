package pguniversal

import (
	"github.com/lukaszkaleta/saas-go/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRating struct {
	db *pg.PgDb
	Id int64
}

func NewPgRating(db *pg.PgDb, id int64) universal.Rating {
	return &PgRating{db, id}
}

func (pgRating *PgRating) Model() *universal.RatingModel {
	return &universal.RatingModel{}
}

func (pgRating *PgRating) Update(newModel *universal.RatingModel) error {
	return nil
}
