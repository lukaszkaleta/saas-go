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
	owner pg.RelationEntity
}

func NewPgMessages(db *pg.PgDb, owner pg.RelationEntity) messages.Messages {
	return &PgMessages{db: db, owner: owner}
}

func (pg *PgMessages) Add(ctx context.Context, recipientId int64, value string) (messages.Message, error) {
	return pg.AddFromModel(ctx, &messages.MessageModel{Value: value, OwnerId: pg.owner.RelationId, RecipientId: recipientId})
}

func (pg *PgMessages) AddGenerated(ctx context.Context, recipientId int64, value string) (messages.Message, error) {
	return pg.AddFromModel(ctx, &messages.MessageModel{ValueGenerated: true, Value: value, OwnerId: pg.owner.RelationId, RecipientId: recipientId})
}

func (pg *PgMessages) AddFromModel(ctx context.Context, model *messages.MessageModel) (messages.Message, error) {
	if model.OwnerId != pg.owner.RelationId {
		return nil, errors.New("Owner inside model and messages does not match")
	}
	messageId := int64(0)
	currentUserId := universal.CurrentUserId(ctx)
	query := fmt.Sprintf("insert into %s (owner_id, user_id, action_created_by_id, value, value_generated) values (@ownerId, @recipientId, @currentUserId, @value, @generated) returning id", pg.owner.TableName)
	row := pg.db.Pool.QueryRow(ctx, query, MessageNamedArgs(model, currentUserId))
	err := row.Scan(&messageId)
	if err != nil {
		return nil, err
	}
	return &PgMessage{
		db:    pg.db,
		Id:    messageId,
		Owner: pg.owner,
	}, nil
}

func (pg *PgMessages) List(ctx context.Context) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId", pg.owner.TableName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.RelationId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapMessage(pg.db, pg.owner))
}

func (pg *PgMessages) ById(ctx context.Context, id int64) (messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where id = @id and owner_id = @ownerId", pg.owner.TableName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id, "ownerId": pg.owner.RelationId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapMessage(pg.db, pg.owner))
}

func (pg *PgMessages) ForRecipient(ctx context.Context, recipient universal.Idable) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId and user_id = @recipientId", pg.owner.TableName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.RelationId, "recipientId": recipient.ID()})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapMessage(pg.db, pg.owner))
}

func (pg *PgMessages) ForRecipientById(ctx context.Context, id int64) ([]messages.Message, error) {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where owner_id = @ownerId and user_id = (select user_id from %s where id = @id)", pg.owner.TableName, pg.owner.TableName)
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"ownerId": pg.owner.RelationId, "id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapMessage(pg.db, pg.owner))
}

func (pg *PgMessages) Acknowledge(ctx context.Context) error {
	currentUserId := universal.CurrentUserId(ctx)
	sql := fmt.Sprintf("update %s set action_read_by_id = @currentUserId where owner_id = @ownerId and action_created_by_id <> @currentUserId", pg.owner.TableName)
	_, err := pg.db.Pool.Exec(ctx, sql, pgx.NamedArgs{"ownerId": pg.owner.RelationId, "currentUserId": currentUserId})
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgMessages) LastQuestions(ctx context.Context) ([]messages.Message, error) {
	sqlTemplate := `
with my_jobs as (
    select
        *,
        rank() over (partition by owner_id order by action_created_at desc)
    from job_message
        where owner_id = @ownerId
)
` + ColumnsSelect() + ` from my_jobs where rank = 1
`
	rows, err := pg.db.Pool.Query(ctx, sqlTemplate, pgx.NamedArgs{"ownerId": pg.owner.RelationId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapMessage(pg.db, pg.owner))
}

func MessageNamedArgs(model *messages.MessageModel, currentUserId *int64) pgx.NamedArgs {
	return pgx.NamedArgs{
		"ownerId":       model.OwnerId,
		"recipientId":   model.RecipientId,
		"currentUserId": currentUserId,
		"value":         model.Value,
		"generated":     model.ValueGenerated,
	}
}
