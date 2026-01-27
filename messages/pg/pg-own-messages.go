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
	db          *pg.PgDb
	owner       pg.TableEntity
	recipientId int64
}

func MyOwnMessages(db *pg.PgDb, owner pg.TableEntity, recipientId int64) messages.Own {
	return PgOwnMessages{db: db, owner: owner, recipientId: recipientId}
}

func (pg PgOwnMessages) LastQuestionsToMe(ctx context.Context) ([]messages.Message, error) {
	currentUserId := universal.CurrentUserId(ctx)
	sqlTemplate := `
		SELECT DISTINCT ON (j.id, jm.action_created_by_id)
			jm.id as id,
			jm.owner_id,
			jm.recipient_id,
		    jm.value,
			jm.action_created_by_id,
			jm.action_created_at,
			jm.action_read_by_id,
			jm.action_read_at
		FROM job j
		INNER JOIN job_message jm
			ON jm.owner_id = j.id
		WHERE
			j.action_created_by_id = @currentUserId and
		  	j.owner_id = @ownerId
		ORDER BY
			j.id,
			jm.action_created_by_id,
			jm.action_created_at DESC;
	`
	query := fmt.Sprintf(sqlTemplate, pg.owner.Name)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"currentUserId": currentUserId, "ownerId": pg.owner.Id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapMessages(rows, pg.db)
}
