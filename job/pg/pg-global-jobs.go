package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgGlobalJobs struct {
	db         *pg.PgDb
	userSearch user.UserSearch
}

func NewPgGlobalJobs(Db *pg.PgDb, userSearch user.UserSearch) job.GlobalJobs {
	return &PgGlobalJobs{db: Db, userSearch: userSearch}
}

func (pgGlobalJobs *PgGlobalJobs) Search(ctx context.Context, input *job.JobSearchInput) ([]*job.JobSearchResult, error) {
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
	return pgGlobalJobs.jobsWithPersons(ctx, rows)
}

func (pgGlobalJobs *PgGlobalJobs) ByQuery(ctx context.Context, query *string) ([]*job.JobSearchResult, error) {
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
	return pgGlobalJobs.jobsWithPersons(ctx, rows)
}

func (pgGlobalJobs *PgGlobalJobs) NearBy(ctx context.Context, radar *universal.RadarModel) ([]*job.JobSearchResult, error) {
	if radar == nil || radar.Position == nil || radar.Position.Lat == 0 || radar.Position.Lon == 0 {
		return pgGlobalJobs.allActiveSearch(ctx)
	}
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
	rows, err := pgGlobalJobs.db.Pool.Query(ctx, sql, pgx.NamedArgs{"lat": radar.Position.Lat, "lon": radar.Position.Lon, "perimeter": radar.Perimeter})
	if err != nil {
		return nil, err
	}
	return pgGlobalJobs.jobsWithPersons(ctx, rows)
}

func (globalJobs *PgGlobalJobs) ActiveById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id and status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapJob(globalJobs.db))
}

func (globalJobs *PgGlobalJobs) ById(ctx context.Context, id int64) (job.Job, error) {
	query := JobSelect() + "where id = @id"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapJob(globalJobs.db))
}

func (globalJobs *PgGlobalJobs) ByIds(ctx context.Context, ids []int64) ([]job.Job, error) {
	query := JobSelect() + "where id = any(@ids)"
	rows, err := globalJobs.db.Pool.Query(ctx, query, pgx.NamedArgs{"ids": ids})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapJob(globalJobs.db))
}

func (globalJobs *PgGlobalJobs) AllActive(ctx context.Context) ([]job.Job, error) {
	query := JobSelect() + " where status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapJob(globalJobs.db))
}

func (globalJobs *PgGlobalJobs) allActiveSearch(ctx context.Context) ([]*job.JobSearchResult, error) {
	query := JobColumnsSelect() + ", 0 as distance, 0 as rank from job where status_published is not null and status_closed is null and status_occupied is null"
	rows, err := globalJobs.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return globalJobs.jobsWithPersons(ctx, rows)
}

func (pgGlobalJobs *PgGlobalJobs) jobsWithPersons(ctx context.Context, rows pgx.Rows) ([]*job.JobSearchResult, error) {
	collectRows, err := pgx.CollectRows(rows, MapSearchJob())
	if err != nil {
		return nil, err
	}

	jobModels := make([]*job.JobModel, len(collectRows))
	for i, r := range collectRows {
		jobModels[i] = r.Model
	}

	userIds := universal.CreatedByIdFromModels(jobModels)
	personModels, err := pgGlobalJobs.userSearch.PersonModelsByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	personMap := make(map[int64]*universal.PersonModel)
	for _, p := range personModels {
		if p != nil {
			personMap[p.Id] = p
		}
	}

	for _, r := range collectRows {
		creatorId := r.Model.Actions.CreatedById()
		if creatorId != nil {
			r.Person = personMap[*creatorId]
		}
	}

	return collectRows, nil
}
