package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgChatOwnerInbox struct {
	db *pg.PgDb
}

func NewPgChatOwnerInbox(db *pg.PgDb) *PgChatOwnerInbox {
	return &PgChatOwnerInbox{db: db}
}

func (p *PgChatOwnerInbox) Last(ctx context.Context) ([]chat.Message, error) {
	currentUserId := universal.CurrentUserId(ctx)

	sqlTemplate := `
WITH latest_messages AS (
    SELECT DISTINCT ON (jm.chat_id)
        jm.*
    FROM job_chat jc
    JOIN job j
        ON j.id = jc.job_id
    JOIN job_message jm
        ON jm.chat_id = jc.id
    WHERE j.action_created_by_id = @currentUserId
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
			"currentUserId": currentUserId,
		},
	)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, MapMessage(p.db))
}

func (p *PgChatOwnerInbox) CountUnread(ctx context.Context) (int, error) {
	currentUserId := universal.CurrentUserId(ctx)

	sqlTemplate := `
SELECT count(*)
FROM job_chat jc
JOIN job j
	ON j.id = jc.job_id
LEFT JOIN job_chat_read jcr
	ON jcr.chat_id = jc.id
	AND jcr.action_updated_by_id = @currentUserId
WHERE j.action_created_by_id = @currentUserId
AND EXISTS (
	SELECT 1
	FROM job_message jm
	WHERE jm.chat_id = jc.id
		AND jm.action_created_by_id <> @currentUserId
		AND jm.id > COALESCE(jcr.last_read_message_id, 0)
)
`

	row := p.db.Pool.QueryRow(
		ctx,
		sqlTemplate,
		pgx.NamedArgs{
			"currentUserId": currentUserId,
		},
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
