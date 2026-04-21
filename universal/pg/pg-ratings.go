package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	err := s.Db.Pool.QueryRow(ctx, query, s.ownerTable.Id).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg == nil {
		return 0, nil
	}
	return int(*avg), nil
}

func (s *PgRatings) AllModels(ctx context.Context) ([]*universal.RatingModel, error) {
	// Map subjectId to the owner table's foreign key column, e.g., job_id
	subjectColumn := s.ownerTable.Name + "_id"
	query := fmt.Sprintf("select id, reviewee_id, score, review_text, review_image_url, action_created_at, action_created_by_id, %s from %s where %s = $1 order by action_created_at desc", subjectColumn, s.ratingTable(), subjectColumn)

	rows, err := s.Db.Pool.Query(ctx, query, s.ownerTable.Id)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, MapRatingModel)
}

type PgRated struct {
	db            *pg.PgDb
	revieweeId    int64
	tableName     string
	subjectColumn string
}

func NewPgRated(db *pg.PgDb, revieweeId int64, tableName string) universal.Rated {
	subjectColumn := tableName[:len(tableName)-len("_rating")] + "_id"
	return &PgRated{
		db:            db,
		revieweeId:    revieweeId,
		tableName:     tableName,
		subjectColumn: subjectColumn,
	}
}

func (s *PgRated) AllModels(ctx context.Context) ([]*universal.RatingModel, error) {
	query := fmt.Sprintf("select id, reviewee_id, score, review_text, review_image_url, action_created_at, action_created_by_id, %s from %s where reviewee_id = $1 order by action_created_at desc", s.subjectColumn, s.tableName)

	rows, err := s.db.Pool.Query(ctx, query, s.revieweeId)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, MapRatingModel)
}

func (s *PgRated) Average(ctx context.Context) (int, error) {
	query := fmt.Sprintf("select avg(score) from %s where reviewee_id = $1", s.tableName)
	var avg *float64
	err := s.db.Pool.QueryRow(ctx, query, s.revieweeId).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg == nil {
		return 0, nil
	}
	return int(*avg), nil
}
