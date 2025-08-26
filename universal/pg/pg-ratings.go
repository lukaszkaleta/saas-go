package postgres

import (
	"naborly/internal/api/rating"postgres2 "naborly/internal/postgres"

)

type PgRatings struct {
	db         *postgres2.PgDb
	ownerTable TableEntity
}

func NewPgRatings(db *postgres2.PgDb, ownerTable TableEntity) rating.Ratings {
	return &PgRatings{db: db, ownerTable: ownerTable}
}

func (s *PgRatings) Add(r *rating.RatingModel) (rating.Rating, error) {
	return rating.NewSolidRating(nil, nil, 0), nil
}

func (s *PgRatings) ById(id int64) (rating.Rating, error) {
	return &rating.SolidRating{}, nil
}
