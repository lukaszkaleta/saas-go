package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
)

type PgOwnerMessages struct {
	db    *pg.PgDb
	jobId int64
}

func NewPgOwnerMessages(db *pg.PgDb, jobId int64) *PgOwnerMessages {
	return &PgOwnerMessages{
		db:    db,
		jobId: jobId,
	}
}

func (p *PgOwnerMessages) Last(ctx context.Context) ([]chat.Message, error) {
	sqlTemplate := `
WITH latest_messages AS (
    SELECT DISTINCT ON (jm.chat_id)
        jm.*
    FROM job_chat jc
    JOIN job_message jm
        ON jm.chat_id = jc.id
    WHERE jc.job_id = @jobId
    ORDER BY jm.chat_id, jm.action_created_at DESC
)
` + ColumnsSelect() + `
FROM latest_messages
ORDER BY action_created_at DESC
`

	rows, err := p.db.Pool.Query(
		ctx,
		sqlTemplate,
		pgx.NamedArgs{
			"jobId": p.jobId,
		},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, MapMessage(p.db))
}
