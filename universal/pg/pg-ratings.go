package pg

import (
	"github.com/lukaszkaleta/saas-go/pg/database"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRatings struct {
	db         *database.PgDb
	ownerTable database.TableEntity
}

func NewPgRatings(db *database.PgDb, ownerTable database.TableEntity) universal.Ratings {
	return &PgRatings{db: db, ownerTable: ownerTable}
}

func (s *PgRatings) Add(r *universal.RatingModel) (universal.Rating, error) {
	return universal.NewSolidRating(nil, nil, 0), nil
}

func (s *PgRatings) ById(id int64) (universal.Rating, error) {
	return &universal.SolidRating{}, nil
}
