package postgres

import (
	"naborly/internal/api/rating"postgres2 "naborly/internal/postgres"

)

type PgRating struct {
	db *postgres2.PgDb
	Id int64
}

func NewPgRating(db *postgres2.PgDb, id int64) rating.Rating {
	return &PgRating{db, id}
}

func (pgRating *PgRating) Model() *rating.RatingModel {
	return &rating.RatingModel{}
}

func (pgRating *PgRating) Update(newModel *rating.RatingModel) error {
	return nil
}
