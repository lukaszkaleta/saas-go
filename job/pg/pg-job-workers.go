package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pguniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgJobWorkers struct {
	db *pg.PgDb
}

func NewPgJobWorkers(db *pg.PgDb) job.JobWorkers {
	return &PgJobWorkers{db: db}
}

func (p *PgJobWorkers) Suggest(ctx context.Context, j job.Job) ([]*universal.PersonModel, error) {
	query := `
		SELECT DISTINCT ` + pguniversal.MapPersonColumnString(func(c string) string { return "u." + c }) + `
		FROM users u
		JOIN task t ON t.user_id = u.id
		JOIN job_category jc_past ON jc_past.job_id = t.job_id
		JOIN job_category jc_current ON jc_current.category_id = jc_past.category_id
		WHERE jc_current.job_id = $1
	`
	rows, err := p.db.Pool.Query(ctx, query, j.ID())
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pguniversal.MapPersonModel)
}
