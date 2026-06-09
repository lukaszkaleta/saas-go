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

func NewPgChats(db *pg.PgDb, owner pg.RelationEntity) chat.ChatsApi {
	return &PgChats{
		db:    db,
		owner: owner,
	}
}

func (c *PgChats) Create(ctx context.Context, workerId int64) (chat.Chat, error) {
	currentUserId := universal.CurrentUserId(ctx)

	query := fmt.Sprintf("insert into job_chat (job_id, worker_id, action_created_by_id) values (@jobId, @workerId, @userId) returning id")
	var id int64
	err := c.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"ownerId":  c.owner.RelationId,
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
