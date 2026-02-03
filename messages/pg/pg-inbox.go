package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgInbox struct {
	db    *pg.PgDb
	owner pg.RelationEntity
}

func NewPgInbox(db *pg.PgDb, owner pg.RelationEntity) messages.Inbox {
	return PgInbox{db: db, owner: owner}
}

func (pg PgInbox) LastQuestions(ctx context.Context) ([]messages.Message, error) {
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
	rows, err := pg.db.Pool.Query(ctx, sqlTemplate, pgx.NamedArgs{"currentUserId": currentUserId})
	if err != nil {
		return nil, err
	}
	return MapMessages(pg.db, pg.owner, rows)
}

func (pg PgInbox) LastAnswers(ctx context.Context) ([]messages.Message, error) {
	currentUserId := universal.CurrentUserId(ctx)
	sqlTemplate := `
with my_tasks as (
    select
        *,
        rank() over (partition by owner_id order by action_created_at desc)
    from job_message
        where recipient_id = @currentUserId
)
` + ColumnsSelect() + ` from my_jobs where rank = 1
`
	rows, err := pg.db.Pool.Query(ctx, sqlTemplate, pgx.NamedArgs{"currentUserId": currentUserId})
	if err != nil {
		return nil, err
	}
	return MapMessages(pg.db, pg.owner, rows)
}
