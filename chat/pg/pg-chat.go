package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgChat struct {
	db *pg.PgDb
	Id int64
}

func (c *PgChat) ID() int64 {
	return c.Id
}

func (c *PgChat) Model(ctx context.Context) (*chat.ChatModel, error) {
	query := "select id, job_id, worker_id, action_created_by_id, action_created_at from job_chat where id=@id"
	rows, err := c.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": c.Id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, MapChatModel)
}

func (c *PgChat) Messages() chat.Messages {
	return NewPgMessages(c.db, c.Id)
}

func MapChatModel(row pgx.CollectableRow) (*chat.ChatModel, error) {
	model := &chat.ChatModel{
		Actions: universal.EmptyActionsModel(),
	}

	actionCreatedModel := universal.EmptyCreatedActionModel()

	err := row.Scan(
		&model.Id,
		&model.JobId,
		&model.WorkerId,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
	)
	if err != nil {
		return nil, err
	}
	model.Actions.List[actionCreatedModel.Name] = actionCreatedModel

	return model, nil
}
