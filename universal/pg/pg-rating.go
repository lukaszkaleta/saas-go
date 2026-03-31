package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRating struct {
	db          *pg.PgDb
	tableEntity pg.TableEntity
}

func (pgRating *PgRating) ID() int64 {
	return pgRating.tableEntity.Id
}

func (pgRating *PgRating) Actions() universal.Actions {
	return NewPgActions(pgRating.db, pgRating.tableEntity)
}

func (pgRating *PgRating) RevieweeId(ctx context.Context) int64 {
	model, err := pgRating.Model(ctx)
	if err != nil {
		return 0
	}
	return model.RevieweeId
}

func (pgRating *PgRating) SubjectId(ctx context.Context) int64 {
	model, err := pgRating.Model(ctx)
	if err != nil {
		return 0
	}
	return model.SubjectId
}

func (pgRating *PgRating) Model(ctx context.Context) (*universal.RatingModel, error) {
	query := fmt.Sprintf("select reviewee_id, score, review_text, review_image_url from %s where id = $1", pgRating.tableEntity.Name)
	var model universal.RatingModel
	model.Review = &universal.DescriptionModel{}

	err := pgRating.db.Pool.QueryRow(ctx, query, pgRating.tableEntity.Id).Scan(
		&model.RevieweeId,
		&model.Score,
		&model.Review.Value,
		&model.Review.ImageUrl,
	)
	if err != nil {
		return universal.EmptyRatingModel(), nil
	}
	model.Id = pgRating.tableEntity.Id

	// Try to get subject_id if it exists, or it might be specific to job_rating
	// In job_rating it is job_id. This is tricky.
	// Based on job_rating: job_id, reviewee_id, score, review_text, review_image_url

	// Re-query to include job_id as subjectId if we are in job_rating
	if pgRating.tableEntity.Name == "job_rating" {
		query = "select job_id, reviewee_id, score, review_text, review_image_url from job_rating where id = $1"
		_ = pgRating.db.Pool.QueryRow(ctx, query, pgRating.tableEntity.Id).Scan(
			&model.SubjectId,
			&model.RevieweeId,
			&model.Score,
			&model.Review.Value,
			&model.Review.ImageUrl,
		)
	}

	actions, _ := pgRating.Actions().Model(ctx)
	model.Actions = actions

	return &model, nil
}

func (pgRating *PgRating) Update(ctx context.Context, rating *universal.RatingModel) error {
	query := fmt.Sprintf("update %s set score = $1, review_text = $2, review_image_url = $3 where id = $4", pgRating.tableEntity.Name)
	_, err := pgRating.db.Pool.Exec(ctx, query, rating.Score, rating.Review.Value, rating.Review.ImageUrl, pgRating.tableEntity.Id)
	return err
}

func NewPgRating(db *pg.PgDb, tableEntity pg.TableEntity) universal.Rating {
	return &PgRating{db, tableEntity}
}

func MapRating(db *pg.PgDb, tableName string) pgx.RowToFunc[universal.Rating] {
	return func(row pgx.CollectableRow) (universal.Rating, error) {
		model, err := MapRatingModel(row)
		if err != nil {
			return nil, err
		}
		pgRating := PgRating{db: db, tableEntity: db.TableEntity(tableName, model.Id)}
		return universal.NewSolidRating(model, &pgRating, model.Id), nil
	}
}

func MapRatingModel(row pgx.CollectableRow) (*universal.RatingModel, error) {
	ratingModel := universal.EmptyRatingModel()
	ratingModel.Actions = universal.EmptyActionsModel()
	created := universal.EmptyCreatedActionModel()
	ratingModel.Actions.List["created"] = created

	columns := row.FieldDescriptions()
	dest := []any{
		&ratingModel.Id,
		&ratingModel.RevieweeId,
		&ratingModel.Score,
		&ratingModel.Review.Value,
		&ratingModel.Review.ImageUrl,
		&created.MadeAt,
		&created.ById,
	}

	if len(columns) > 7 {
		dest = append(dest, &ratingModel.SubjectId)
	}

	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	return ratingModel, nil
}
