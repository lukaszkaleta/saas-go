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

func (pgGlobalJobs *PgGlobalJobs) Search(ctx context.Context, input *job.JobSearchInput) ([]*job.JobSearchOutput, error) {
	if len(*input.Query) <= 2 {
		return pgGlobalJobs.NearBy(ctx, input.Radar)
	}
	if input.Radar == nil {
		return pgGlobalJobs.ByQuery(ctx, input.Query)
	}

	// Combine NearBy and ByQuery with window function
	// First create ranked window query
	ftsSql := `
		WITH fts_limited AS (
		  	` + JobColumnsSelect() + `,
			earth_point,
			ts_rank_cd(search_vector, q) AS rank 
		  FROM job,
			websearch_to_tsquery('norwegian', @query) q
		  WHERE 
			search_vector @@ q and
			status_published is not null and 
			status_closed is null and 
			status_occupied is null
		  ORDER BY rank DESC
		  LIMIT 2000
		)
	`
	// Then near by and order by distance
	jobSql := JobColumnsSelectWithPrefix("p") + `,
          earth_distance(p.earth_point, ll_to_earth(@lat, @lon)) AS distance,
		  p.rank as rank
		FROM fts_limited p
		WHERE 
			p.earth_point <@ earth_box(ll_to_earth(@lat, @lon), @perimeter) and
			p.status_published is not null and 
			p.status_closed is null and 
			p.status_occupied is null
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
	return MapSearchJobs(rows)
}

func (pgGlobalJobs *PgGlobalJobs) ByQuery(ctx context.Context, query *string) ([]*job.JobSearchOutput, error) {
	sql := JobColumnsSelect() + `,
			0 as distance,
			ts_rank_cd(search_vector, q) AS rank
		  FROM job,
			websearch_to_tsquery('norwegian', @query) q
		  WHERE 
			search_vector @@ q and
			status_published is not null and 
			status_closed is null and 
			status_occupied is null
		  ORDER BY rank DESC
		  LIMIT 2000
`
	rows, err := pgGlobalJobs.db.Pool.Query(ctx, sql, pgx.NamedArgs{"query": query})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapSearchJobs(rows)
}

func (globalJobs *PgGlobalJobs) NearBy(ctx context.Context, radar *universal.RadarModel) ([]*job.JobSearchOutput, error) {
	sql := JobColumnsSelect() + `,
			earth_distance(earth_point, ll_to_earth(@lat, @lon)) AS distance,
			earth_distance(earth_point, ll_to_earth(@lat, @lon)) AS rank
		from job
		where 
			earth_point <@ earth_box(ll_to_earth(@lat, @lon), @perimeter) and
			status_published is not null and 
			status_closed is null and 
			status_occupied is null
`
	rows, err := globalJobs.db.Pool.Query(ctx, sql, pgx.NamedArgs{"lat": radar.Position.Lat, "lon": radar.Position.Lon, "perimeter": radar.Perimeter})
	if err != nil {
		return nil, err
	}
	return MapSearchJobs(rows)
}

func (globalJobs *PgGlobalJobs) ActiveById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id and status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return MapJob(globalJobs.db)(rows)
	}
	return nil, nil
}

func (globalJobs *PgGlobalJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return MapJob(globalJobs.db)(rows)
	}
	return nil, nil
}

func (globalJobs *PgGlobalJobs) ByIds(ctx context.Context, ids []int64) ([]job.Job, error) {
	query := JobSelect() + "where id = any(@ids)"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"ids": ids})
	if err != nil {
		return nil, err
	}
	return MapJobs(globalJobs.db, rows)
}

func (globalJobs *PgGlobalJobs) AllActive(ctx context.Context) ([]job.Job, error) {
	query := JobSelect() + " where status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return MapJobs(globalJobs.db, rows)
}
