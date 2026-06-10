package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessage struct {
	db *pg.PgDb
	Id int64
}

func NewPgMessage(db *pg.PgDb, id int64) *PgMessage {
	return &PgMessage{db: db, Id: id}
}

func (m *PgMessage) ID() int64 {
	return m.Id
}

func (m *PgMessage) Acknowledge(ctx context.Context) error {
	currentUserId := universal.CurrentUserId(ctx)
	_, err := m.db.Pool.Exec(ctx, `
		INSERT INTO job_chat_read (
			chat_id,
			last_read_message_id,
			action_updated_by_id,
			action_updated_at
		)
		SELECT
			jm.chat_id,
			jm.id,
			@currentUserId,
			now()
		FROM job_message jm
		WHERE jm.id = @messageId
		ON CONFLICT (chat_id, action_updated_by_id)
		DO UPDATE
		SET
			last_read_message_id = GREATEST(
				job_chat_read.last_read_message_id,
				EXCLUDED.last_read_message_id
			),
			action_updated_at = now()
	`, pgx.NamedArgs{
		"messageId":     m.Id,
		"currentUserId": currentUserId,
	})

	return err
}

func (m *PgMessage) Model(ctx context.Context) (*chat.MessageModel, error) {
	query := ColumnsSelect() + " from job_message where id=@id"
	rows, err := m.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": m.Id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, MapMessageModel)
}

func MapMessageModel(row pgx.CollectableRow) (*chat.MessageModel, error) {
	model := &chat.MessageModel{
		Actions: universal.EmptyActionsModel(),
	}

	actionCreatedModel := universal.EmptyCreatedActionModel()

	err := row.Scan(
		&model.Id,
		&model.ChatId,
		&model.Value,
		&model.ValueGenerated,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
	)
	if err != nil {
		return nil, err
	}
	model.Actions.List[actionCreatedModel.Name] = actionCreatedModel

	return model, nil
}

func Columns() []string {
	return []string{
		"id",
		"chat_id",
		"value",
		"value_generated",
		"action_created_by_id",
		"action_created_at",
	}
}

func ColumnString() string {
	return strings.Join(Columns(), ",")
}

func ColumnsSelect() string {
	return "select " + ColumnString()
}

func MapMessage(db *pg.PgDb) pgx.RowToFunc[chat.Message] {
	return func(row pgx.CollectableRow) (chat.Message, error) {
		model, err := MapMessageModel(row)
		if err != nil {
			return nil, err
		}
		pgMessage := NewPgMessage(db, model.Id)
		return chat.NewSolidMessage(model, pgMessage), nil
	}

}
