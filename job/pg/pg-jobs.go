package pgjob

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgJobs struct {
	db  *pg.PgDb
	Ids []int
}

func (pgJobs *PgJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = $1"
	rows, err := pgJobs.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return MapJob(pgJobs.db)(rows)
}

func (pgJobs *PgJobs) Add(ctx context.Context, model *job.JobModel) (job.Job, error) {
	jobId := int64(0)
	query := "INSERT INTO job (description_value, description_image_url, position_latitude, position_longitude, address_line_1, address_line_2, address_city, address_postal_code, address_district, price_value, price_currency, rating, tags, action_created_by_id) VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) returning id"
	currentUser := user.CurrentUser(ctx)
	row := pgJobs.db.Pool.QueryRow(
		ctx,
		query,
		model.Description.Value,
		model.Description.ImageUrl,
		model.Position.Lat,
		model.Position.Lon,
		model.Address.Line1,
		model.Address.Line2,
		model.Address.City,
		model.Address.PostalCode,
		model.Address.District,
		model.Price.Value,
		model.Price.Currency,
		model.Rating,
		model.Tags,
		currentUser.Id,
	)
	err := row.Scan(&jobId)
	if err != nil {
		return nil, err
	}
	pgJob := PgJob{
		db: pgJobs.db,
		Id: jobId,
	}
	return job.NewSolidJob(
		&job.JobModel{
			Id:          jobId,
			Description: &universal.DescriptionModel{},
			Position:    model.Position,
			Address:     model.Address,
			Price:       &universal.PriceModel{},
			Rating:      model.Rating,
			State:       job.JobStatus{Draft: time.Now()},
		},
		&pgJob), nil
}

func (pgJobs *PgJobs) List(ctx context.Context) ([]job.Job, error) {
	return nil, errors.New("All jobs can not be listed")
}

func MapJobs(db *pg.PgDb, rows pgx.Rows) ([]job.Job, error) {
	jobs := []job.Job{}
	for rows.Next() {
		mJob, err := MapJob(db)(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, mJob)
	}
	return jobs, nil
}

// Relation

type PgRelationJobs struct {
	Db       *pg.PgDb
	Jobs     *PgJobs
	Relation pg.RelationEntity
}

func NewPgRelationJobs(pfJobs *PgJobs, relation pg.RelationEntity) PgRelationJobs {
	return PgRelationJobs{
		Db:       pfJobs.db,
		Jobs:     pfJobs,
		Relation: relation,
	}
}

func (p PgRelationJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	return p.Jobs.ById(ctx, id)
}

func (p PgRelationJobs) Add(ctx context.Context, jobModel *job.JobModel) (job.Job, error) {
	newJob, err := p.Jobs.Add(ctx, jobModel)
	if err != nil {
		return newJob, err
	}
	query := fmt.Sprintf("INSERT INTO %s(job_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err = p.Db.Pool.Exec(ctx, query, newJob.Model().Id, p.Relation.RelationId)
	if err != nil {
		return newJob, err
	}
	return newJob, nil
}

func (p PgRelationJobs) Join(ctx context.Context, jobId int64) error {
	query := fmt.Sprintf("INSERT INTO %s(job_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err := p.Db.Pool.Exec(ctx, query, jobId, p.Relation.RelationId)
	if err != nil {
		return err
	}
	return nil
}

func (p PgRelationJobs) List(ctx context.Context) ([]job.Job, error) {
	query := fmt.Sprintf("%s where id in (select job_id from %s where %s = $1)", JobSelect(), p.Relation.TableName, p.Relation.ColumnName)
	rows, err := p.Db.Pool.Query(ctx, query, p.Relation.RelationId)
	if err != nil {
		return nil, err
	}
	return MapJobs(p.Db, rows)
}
