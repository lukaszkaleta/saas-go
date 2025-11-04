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
)

type PgJobs struct {
	Db  *pg.PgDb
	Ids []int
}

func (pgJobs *PgJobs) Add(model *job.JobModel) (job.Job, error) {
	jobId := int64(0)
	query := "INSERT INTO job (description_value, description_image_url, position_latitude, position_longitude, address_line_1, address_line_2, address_city, address_postal_code, address_district, price_value, price_currency) VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 ) returning id"
	row := pgJobs.Db.Pool.QueryRow(
		context.Background(),
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
	)
	err := row.Scan(&jobId)
	if err != nil {
		return nil, err
	}
	pgJob := PgJob{
		Db: pgJobs.Db,
		Id: jobId,
	}
	return job.NewSolidJob(
		&job.JobModel{
			Id:          jobId,
			Description: &universal.DescriptionModel{},
			Position:    model.Position,
			Address:     model.Address,
			Price:       &universal.PriceModel{},
			State:       job.JobStatus{Draft: time.Now()},
		},
		&pgJob,
		jobId,
	), nil
}

func (pgJobs *PgJobs) List() ([]job.Job, error) {
	return nil, errors.New("All jobs can not be listed")
}

func MapJobs(rows pgx.Rows, db *pg.PgDb) ([]job.Job, error) {
	jobs := []job.Job{}
	id := int64(0)
	for rows.Next() {
		pgJob := &PgJob{Db: db, Id: id}
		jobModel, err := MapJob(rows)
		if err != nil {
			return nil, err
		}
		solidJob := job.NewSolidJob(
			jobModel,
			pgJob,
			id)
		jobs = append(jobs, solidJob)
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
		Db:       pfJobs.Db,
		Jobs:     pfJobs,
		Relation: relation,
	}
}

func (p PgRelationJobs) Add(jobModel *job.JobModel) (job.Job, error) {
	newJob, err := p.Jobs.Add(jobModel)
	if err != nil {
		return newJob, err
	}
	query := fmt.Sprintf("INSERT INTO %s(job_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err = p.Db.Pool.Exec(context.Background(), query, newJob.Model().Id, p.Relation.RelationId)
	if err != nil {
		return newJob, err
	}
	return newJob, nil
}

func (p PgRelationJobs) Join(jobId int64) error {
	query := fmt.Sprintf("INSERT INTO %s(job_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err := p.Db.Pool.Exec(context.Background(), query, jobId, p.Relation.RelationId)
	if err != nil {
		return err
	}
	return nil
}

func (p PgRelationJobs) List() ([]job.Job, error) {
	query := fmt.Sprintf("select * from job where id in (select job_id from %s where %s = $1)", p.Relation.TableName, p.Relation.ColumnName)
	rows, err := p.Db.Pool.Query(context.Background(), query, p.Relation.RelationId)
	if err != nil {
		return nil, err
	}
	return MapJobs(rows, p.Db)
}
