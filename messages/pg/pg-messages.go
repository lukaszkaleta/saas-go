package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessages struct {
	Db      *pg.PgDb
	OwnerId int64
}

func NewPgMessages(db *pg.PgDb, ownerId int64) messages.Messages {
	return &PgMessages{Db: db, OwnerId: ownerId}
}

func (pg *PgMessages) Add(ctx context.Context, value string) (messages.Message, error) {
	return pg.AddFromModel(ctx, &messages.MessageModel{Value: value, OwnerId: pg.OwnerId})
}

func (pg *PgMessages) AddFromModel(ctx context.Context, model *messages.MessageModel) (messages.Message, error) {
	if model.OwnerId != pg.OwnerId {
		return nil, errors.New("owner inside model and messages does not match")
	}
	messageId := int64(0)
	currentUserId := universal.CurrentUserId(ctx)
	query := "insert into message (owner_id, action_created_by_id, value) values (@ownerId, @currentUserId, @value) returning id"
	row := pg.Db.Pool.QueryRow(ctx, query, MessageNamedArgs(model, currentUserId))
	err := row.Scan(&messageId)
	if err != nil {
		return nil, err
	}
	return &PgMessage{
		Db:      pg.Db,
		Id:      messageId,
		OwnerId: model.OwnerId,
	}, nil
}

func (pg *PgMessages) List(ctx context.Context) ([]messages.Message, error) {
	query := "select id, owner_id, value, action_created_by_id, action_created_at from message where owner_id = @ownerId"
	rows, err := pg.Db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.OwnerId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapMessages(rows, pg.Db)
}

func MessageNamedArgs(model *messages.MessageModel, currentUserId *int64) pgx.NamedArgs {
	return pgx.NamedArgs{
		"ownerId":       model.OwnerId,
		"currentUserId": currentUserId,
		"value":         model.Value,
	}
}

func MapMessages(rows pgx.Rows, db *pg.PgDb) ([]messages.Message, error) {
	msgs := []messages.Message{}
	id := int64(0)
	for rows.Next() {
		pgMessage := &PgMessage{Db: db, Id: id}
		msgModel, err := MapMessageModel(rows)
		if err != nil {
			return nil, err
		}
		solidMessage := messages.NewSolidMessage(
			msgModel,
			pgMessage,
			id)
		msgs = append(msgs, solidMessage)
	}
	return msgs, nil
}
