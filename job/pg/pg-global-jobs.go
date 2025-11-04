package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgGlobalJobs struct {
	Db *pg.PgDb
}

func NewPgGlobalJobs(Db *pg.PgDb) job.GlobalJobs {
	return &PgGlobalJobs{Db}
}

func (globalJobs *PgGlobalJobs) NearBy(radar *universal.RadarModel) ([]job.Job, error) {
	id := int64(0)
	query := "select * from job where status_published is not null and status_closed is null and status_occupied is null"
	jobs := []job.Job{}
	rows, err := globalJobs.Db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pgJob := &PgJob{Db: globalJobs.Db, Id: id}
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

func (globalJobs *PgGlobalJobs) ById(id int64) (job.Job, error) {
	query := "select * from job where id = @id and status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.Db.Pool.Query(context.Background(), query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}

	pgJob := &PgJob{Db: globalJobs.Db, Id: id}
	jobModel, err := MapJob(rows)
	if err != nil {
		return nil, err
	}

	return job.NewSolidJob(jobModel, pgJob, jobModel.Id), nil
}
