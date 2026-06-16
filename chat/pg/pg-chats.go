package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgChats struct {
	db    *pg.PgDb
	owner pg.RelationEntity
}

func (c *PgChats) LastMessages(ctx context.Context) ([]chat.Message, error) {
	query := `
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

	rows, err := c.db.Pool.Query(ctx, query, pgx.NamedArgs{
		"jobId": c.owner.RelationId,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, MapMessage(c.db))
}

func (c *PgChats) Delete(ctx context.Context) error {
	query := fmt.Sprintf("delete from %s where %s = @relationId", c.owner.TableName, c.owner.ColumnName)
	_, err := c.db.Pool.Exec(ctx, query, pgx.NamedArgs{
		"relationId": c.owner.RelationId,
	})
	return err
}

func (c *PgChats) ById(ctx context.Context, id int64) (chat.Chat, error) {
	pgChat := &PgChat{
		db: c.db,
		Id: id,
	}

	model, err := pgChat.Model(ctx)
	if err != nil {
		return nil, err
	}

	return chat.NewSolidChat(model, pgChat, id, pgChat.Messages()), nil
}

func (c *PgChats) ByWorkerId(ctx context.Context, id int64) (chat.Chat, error) {
	chatId, err := c.findChatId(ctx, id)
	if err != nil {
		return nil, err
	}

	if chatId == 0 {
		return nil, pgx.ErrNoRows
	}

	return c.ById(ctx, chatId)
}

func NewPgChats(db *pg.PgDb, owner pg.RelationEntity) chat.Chats {
	return &PgChats{
		db:    db,
		owner: owner,
	}
}

func (c *PgChats) findChatId(ctx context.Context, workerId int64) (int64, error) {
	var id int64
	checkQuery := "select id from job_chat where job_id = @jobId and worker_id = @workerId"
	err := c.db.Pool.QueryRow(ctx, checkQuery, pgx.NamedArgs{
		"jobId":    c.owner.RelationId,
		"workerId": workerId,
	}).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (c *PgChats) Ensure(ctx context.Context, workerId int64) (chat.Chat, error) {
	currentUserId := universal.CurrentUserId(ctx)

	id, err := c.findChatId(ctx, workerId)
	if err != nil {
		return nil, err
	}

	if id != 0 {
		pgChat := &PgChat{
			db: c.db,
			Id: id,
		}

		model, err := pgChat.Model(ctx)
		if err != nil {
			return nil, err
		}

		return chat.NewSolidChat(model, pgChat, id, pgChat.Messages()), nil
	}

	query := fmt.Sprintf("insert into job_chat (job_id, worker_id, action_created_by_id) values (@jobId, @workerId, @userId) returning id")
	err = c.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"jobId":    c.owner.RelationId,
		"workerId": workerId,
		"userId":   currentUserId,
	}).Scan(&id)
	if err != nil {
		return nil, err
	}

	pgChat := &PgChat{
		db: c.db,
		Id: id,
	}

	model, err := pgChat.Model(ctx)
	if err != nil {
		return nil, err
	}

	return chat.NewSolidChat(model, pgChat, id, pgChat.Messages()), nil
}
