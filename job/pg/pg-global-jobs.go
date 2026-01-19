package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgGlobalJobs struct {
	db *pg.PgDb
}

func NewPgGlobalJobs(Db *pg.PgDb) job.GlobalJobs {
	return &PgGlobalJobs{Db}
}

func (pgGlobalJobs *PgGlobalJobs) ByQuery(ctx context.Context, query string) ([]job.Job, error) {
	sql := JobSelect() + `, to_tsquery('norwegian', $1) query
		WHERE search_vector @@ query
		ORDER BY ts_rank(search_vector, query) DESC;
	`
	rows, err := pgGlobalJobs.db.Pool.Query(ctx, sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapJobs(pgGlobalJobs.db, rows)
}

func (globalJobs *PgGlobalJobs) NearBy(ctx context.Context, radar *universal.RadarModel) ([]job.Job, error) {
	query := JobSelect() + " where status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return MapJobs(globalJobs.db, rows)
}

func (globalJobs *PgGlobalJobs) ActiveById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id and status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	return MapJob(globalJobs.db)(rows)
}

func (globalJobs *PgGlobalJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	return MapJob(globalJobs.db)(rows)
}
