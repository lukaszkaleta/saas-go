package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessages struct {
	db     *pg.PgDb
	chatId int64
}

func (m *PgMessages) AddGenerated(ctx context.Context, in string) (chat.Message, error) {
	return m.insert(ctx, in, true)
}

func (m *PgMessages) Create(ctx context.Context, in string) (chat.Message, error) {
	return m.insert(ctx, in, false)
}

func (m *PgMessages) ById(ctx context.Context, id int64) (chat.Message, error) {
	query := ColumnsSelect() + " from job_message where id=@id and chat_id=@chatId"
	rows, err := m.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id, "chatId": m.chatId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, MapMessage(m.db))
}

func NewPgMessages(db *pg.PgDb, chatId int64) chat.Messages {
	return &PgMessages{
		db:     db,
		chatId: chatId,
	}
}

func (m *PgMessages) List(ctx context.Context) ([]chat.Message, error) {
	query := ColumnsSelect() + " from job_message where chat_id=@chatId"
	rows, err := m.db.Pool.Query(ctx, query, pgx.NamedArgs{"chatId": m.chatId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, MapMessage(m.db))
}

func (m *PgMessages) Acknowledge(ctx context.Context) error {
	currentUserId := universal.CurrentUserId(ctx)
	query := "update job_message set action_read_at = now(), action_read_by_id = @userId where chat_id = @chatId"
	_, err := m.db.Pool.Exec(ctx, query, pgx.NamedArgs{"userId": currentUserId, "chatId": m.chatId})
	return err
}

func (m *PgMessages) insert(ctx context.Context, in string, generated bool) (chat.Message, error) {
	query := "insert into job_message (chat_id, value, value_generated, action_created_by_id, action_created_at) values (@chatId, @value, @valueGenerated, @byId, @at) returning id"
	args := pgx.NamedArgs{
		"chatId":         m.chatId,
		"value":          in,
		"valueGenerated": generated,
		"byId":           universal.CurrentUserId(ctx),
		"at":             time.Now(),
	}

	var id int64
	err := m.db.Pool.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return nil, err
	}

	return chat.NewSolidMessage(&chat.MessageModel{
		Id:             id,
		ChatId:         m.chatId,
		Value:          in,
		ValueGenerated: false,
	}, NewPgMessage(m.db, id)), nil
}
