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
SELECT DISTINCT ON (o.id, m.action_created_by_id)
	m.id, 
  	m.owner_id, 
  	m.recipient_id, 
  	m.value, 
  	m.action_created_by_id, 
  	m.action_created_at, 
  	m.action_read_by_id, action_read_at
FROM %s o
INNER JOIN %s_message m
    ON m.owner_id = o.id
WHERE
    o.action_created_by_id = @currentUserId
ORDER BY
    o.id,
    m.action_created_by_id,
    m.action_created_at DESC;
	`
	query := fmt.Sprintf(sqlTemplate, pg.ownerName, pg.ownerName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"currentUserId": currentUserId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapMessages(pg.db, rows)
}
