package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgOwnMessages struct {
	db        *pg.PgDb
	ownerName string
}

func MyOwnMessages(db *pg.PgDb, ownerName string) messages.Own {
	return PgOwnMessages{db: db, ownerName: ownerName}
}

func (pg PgOwnMessages) LastQuestionsToMe(ctx context.Context) ([]messages.Message, error) {
	currentUserId := universal.CurrentUserId(ctx)
	sqlTemplate := `
with my_jobs as (
    select
        *,
        rank() over (partition by owner_id order by action_created_at desc)
    from job_message
        where owner_id in (select id from job where job.action_created_by_id = @currentUserId)
)
` + ColumnsSelect() + ` from my_jobs where rank = 1
`
	query := fmt.Sprintf(sqlTemplate, pg.ownerName, pg.ownerName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"currentUserId": currentUserId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapMessages(pg.db, rows)
}
