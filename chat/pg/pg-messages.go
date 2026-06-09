package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/chat"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
)

type PgMessages struct {
	db     *pg.PgDb
	chatId int64
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
