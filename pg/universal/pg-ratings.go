package universal

import (
	"github.com/lukaszkaleta/saas-go/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRatings struct {
	db         *pg.PgDb
	ownerTable pg.TableEntity
}

func NewPgRatings(db *pg.PgDb, ownerTable pg.TableEntity) universal.Ratings {
	return &PgRatings{db: db, ownerTable: ownerTable}
}

func (s *PgRatings) Add(r *universal.RatingModel) (universal.Rating, error) {
	return universal.NewSolidRating(nil, nil, 0), nil
}

func (s *PgRatings) ById(id int64) (universal.Rating, error) {
	return &universal.SolidRating{}, nil
}
