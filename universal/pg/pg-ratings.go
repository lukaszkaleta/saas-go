package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRatings struct {
	Db         *pg.PgDb
	ownerTable pg.TableEntity
}

func NewPgRatings(db *pg.PgDb, ownerTable pg.TableEntity) universal.Ratings {
	return &PgRatings{Db: db, ownerTable: ownerTable}
}

func (s *PgRatings) ratingTable() string {
	return s.ownerTable.Name + "_rating"
}

func (s *PgRatings) Add(ctx context.Context, r *universal.RatingModel) (universal.Rating, error) {
	// Map subjectId to the owner table's foreign key column, e.g., job_id
	subjectColumn := s.ownerTable.Name + "_id"
	query := fmt.Sprintf("insert into %s (%s, reviewee_id, score, review_text, review_image_url, action_created_by_id) values ($1, $2, $3, $4, $5, $6) returning id", s.ratingTable(), subjectColumn)

	createdById := universal.CurrentUserId(ctx)
	var id int64
	err := s.Db.Pool.QueryRow(ctx, query, r.SubjectId, r.RevieweeId, r.Score, r.Review.Value, r.Review.ImageUrl, createdById).Scan(&id)
	if err != nil {
		return nil, err
	}

	return NewPgRating(s.Db, pg.TableEntity{Name: s.ratingTable(), Id: id}), nil
}

func (s *PgRatings) ById(ctx context.Context, id int64) (universal.Rating, error) {
	return NewPgRating(s.Db, pg.TableEntity{Name: s.ratingTable(), Id: id}), nil
}

func (s *PgRatings) Average(ctx context.Context) (int, error) {
	query := fmt.Sprintf("select avg(score) from %s where %s_id = $1", s.ratingTable(), s.ownerTable.Name)
	var avg *float64
	err := s.Db.Pool.QueryRow(context.Background(), query, s.ownerTable.Id).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg == nil {
		return 0, nil
	}
	return int(*avg), nil
}
