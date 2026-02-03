package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessages struct {
	db    *pg.PgDb
	owner pg.TableEntity
}

func NewPgMessages(db *pg.PgDb, owner pg.TableEntity) messages.Messages {
	return &PgMessages{db: db, owner: owner}
}

func (pg *PgMessages) Add(ctx context.Context, recipientId int64, value string) (messages.Message, error) {
	return pg.AddFromModel(ctx, &messages.MessageModel{Value: value, OwnerId: pg.owner.Id, RecipientId: recipientId})
}

func (pg *PgMessages) AddFromModel(ctx context.Context, model *messages.MessageModel) (messages.Message, error) {
	if model.OwnerId != pg.owner.Id {
		return nil, errors.New("owner inside model and messages does not match")
	}
	messageId := int64(0)
	currentUserId := universal.CurrentUserId(ctx)
	query := fmt.Sprintf("insert into %s (owner_id, recipient_id, action_created_by_id, value) values (@ownerId, @recipientId, @currentUserId, @value) returning id", pg.owner.Name)
	row := pg.db.Pool.QueryRow(ctx, query, MessageNamedArgs(model, currentUserId))
	err := row.Scan(&messageId)
	if err != nil {
		return nil, err
	}
	return &PgMessage{
		Db:      pg.db,
		Id:      messageId,
		OwnerId: model.OwnerId,
	}, nil
}

func (pg *PgMessages) List(ctx context.Context) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId", pg.owner.Name)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.Id})
	if err != nil {
		return nil, err
	}
	return MapMessages(pg.db, rows)
}

func (pg *PgMessages) ById(ctx context.Context, id int64) (messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where id = @id and owner_id = @ownerId", pg.owner.Name)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id, "ownerId": pg.owner.Id})
	if err != nil {
		return nil, err
	}
	mapMessages, err := MapMessages(pg.db, rows)
	if err != nil {
		return nil, err
	}
	return mapMessages[0], nil
}

func (pg *PgMessages) ForRecipient(ctx context.Context, recipient universal.Idable) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId and recipient_id = @recipientId", pg.owner.Name)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.Id, "recipientId": recipient.ID()})
	if err != nil {
		return nil, err
	}
	return MapMessages(pg.db, rows)
}

func (pg *PgMessages) ForRecipientById(ctx context.Context, id int64) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId and recipient_id = (select recipient_id from %s where id = @id)", pg.owner.Name, pg.owner.Name)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.Id, "id": id})
	if err != nil {
		return nil, err
	}
	return MapMessages(pg.db, rows)
}

func MessageNamedArgs(model *messages.MessageModel, currentUserId *int64) pgx.NamedArgs {
	return pgx.NamedArgs{
		"ownerId":       model.OwnerId,
		"recipientId":   model.RecipientId,
		"currentUserId": currentUserId,
		"value":         model.Value,
	}
}

func MapMessages(db *pg.PgDb, rows pgx.Rows) ([]messages.Message, error) {
	msgs := []messages.Message{}
	defer rows.Close()
	for rows.Next() {
		msg, err := MapMessage(db, rows)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
