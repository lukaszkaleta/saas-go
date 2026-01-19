package pgjob

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
)

type PgJobSearch struct {
	db *pg.PgDb
}

func NewPgJobSearch(db *pg.PgDb) *PgJobSearch {
	return &PgJobSearch{db: db}
}

func (pgJobSearch *PgJobSearch) ByQuery(ctx context.Context, query string) ([]job.Job, error) {
	sql := `
		SELECT *, ts_rank(search_vector, query) as rank
		FROM job, to_tsquery('norwegian', $1) query
		WHERE search_vector @@ query
		ORDER BY rank DESC;
	`
	rows, err := pgJobSearch.db.Pool.Query(ctx, sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapJobs(pgJobSearch.db, rows)
}
