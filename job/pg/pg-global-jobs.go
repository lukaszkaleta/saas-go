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
	return &PgGlobalJobs{db: Db}
}

func (pgGlobalJobs *PgGlobalJobs) Search(ctx context.Context, input job.JobSearchInput) ([]job.Job, error) {
	if true {
		return pgGlobalJobs.AllActive(ctx)
	}

	if len(input.Query) <= 2 {
		return pgGlobalJobs.NearBy(ctx, input.Radar)
	}
	if input.Radar == nil {
		return pgGlobalJobs.ByQuery(ctx, input.Query)
	}

	ftsSql := `
		WITH fts_limited AS (
		  SELECT
			id,
			earth_point,
			ts_rank_cd(search_vector, q) AS rank
		  FROM job,
			   websearch_to_tsquery('norwegian', $1) q
		  WHERE search_vector @@ q
		  ORDER BY rank DESC
		  LIMIT 2000
		)
	`
	jobSql := JobColumnsSelect() + `
		  p.rank,
		  earth_distance(p.earth_point, ll_to_earth(@lat, @lon)) AS distance
		FROM fts_limited p
		WHERE p.earth_point
			  <@ earth_box(ll_to_earth(@lat, @lon), $perimeter)
		ORDER BY
		  p.rank DESC,
		  distance ASC
		LIMIT 200;
	`
	sql := ftsSql + " " + jobSql
	args := pgx.NamedArgs{
		"query":     input.Query,
		"lat":       input.Radar.Position.Lat,
		"lon":       input.Radar.Position.Lon,
		"perimeter": input.Radar.Perimeter,
	}
	rows, err := pgGlobalJobs.db.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//return MapJobsWith(rows, MapSearchJob(pgGlobalJobs.db))
	return nil, nil
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
	return MapJob(globalJobs.db)(rows)
}

func (globalJobs *PgGlobalJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	return MapJob(globalJobs.db)(rows)
}

func (globalJobs *PgGlobalJobs) AllActive(ctx context.Context) ([]job.Job, error) {
	query := JobSelect() + " where status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return MapJobs(globalJobs.db, rows)
}
