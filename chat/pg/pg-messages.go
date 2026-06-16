package pg

import (
	"context"
	"errors"
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
	query := `
		insert into job_chat_read (chat_id, last_read_message_id, action_updated_by_id, action_updated_at)
		select @chatId, max(id), @userId, now()
		from job_message
		where chat_id = @chatId
		on conflict (chat_id, action_updated_by_id)
		do update set last_read_message_id = GREATEST(job_chat_read.last_read_message_id, excluded.last_read_message_id), action_updated_at = excluded.action_updated_at
	`
	_, err := m.db.Pool.Exec(ctx, query, pgx.NamedArgs{"userId": currentUserId, "chatId": m.chatId})
	return err
}

func (m *PgMessages) AddGenerated(ctx context.Context, value string) (chat.Message, error) {
	return m.insert(ctx, value, true)
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
		ValueGenerated: generated,
	}, NewPgMessage(m.db, id)), nil
}

func (m *PgMessages) LastReadMessageId(ctx context.Context) (int64, error) {
	currentUserId := universal.CurrentUserId(ctx)
	query := "select last_read_message_id from job_chat_read where chat_id = @chatId and action_updated_by_id = @userId"
	var id int64
	err := m.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"chatId": m.chatId, "userId": currentUserId}).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (m *PgMessages) LastReadMessageAt(ctx context.Context) (time.Time, error) {
	currentUserId := universal.CurrentUserId(ctx)
	query := "select action_updated_at from job_chat_read where chat_id = @chatId and action_updated_by_id = @userId"
	var at time.Time
	err := m.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"chatId": m.chatId, "userId": currentUserId}).Scan(&at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return at, nil
}
