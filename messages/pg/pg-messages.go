package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessages struct {
	Db      *pg.PgDb
	Context context.Context
	Ids     []int64
}

func NewPgMessages(db *pg.PgDb) *PgMessages {
	return &PgMessages{Db: db}
}

func (pg *PgMessages) Add(ctx context.Context, model *messages.Model) (messages.Message, error) {
	messageId := int64(0)
	currentUserId := universal.CurrentUserId(ctx)
	query := "insert into messages (owner_id, action_created_by_id, value) values (@owner_id, @currentUserId, @value) returning id"
	row := pg.Db.Pool.QueryRow(ctx, query, MessageNamedArgs(model, currentUserId))
	err := row.Scan(&messageId)
	if err != nil {
		return nil, err
	}
	return PgMessage{
		Db:      pg.Db,
		Id:      messageId,
		OwnerId: model.OwnerId,
	}, nil
}

func MessageNamedArgs(model *messages.Model, currentUserId *int64) pgx.NamedArgs {
	return pgx.NamedArgs{
		"ownerId":       model.OwnerId,
		"currentUserId": currentUserId,
		"value":         model.Value,
	}
}
