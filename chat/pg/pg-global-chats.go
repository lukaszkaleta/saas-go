package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	pg "github.com/lukaszkaleta/saas-go/database/pg"
)

type PgGlobalChats struct {
	db *pg.PgDb
}

func NewPgGlobalChats(db *pg.PgDb) *PgGlobalChats {
	return &PgGlobalChats{db: db}
}

func (p *PgGlobalChats) RelationIds(ctx context.Context, chatIds []int64) (map[int64]int64, error) {
	query := `
		SELECT jc.id, jc.job_id
		FROM job_chat jc
		WHERE jc.id = ANY(@chatIds)
	`
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"chatIds": chatIds})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int64)
	for rows.Next() {
		var chatId, ownerId int64
		if err := rows.Scan(&chatId, &ownerId); err != nil {
			return nil, err
		}
		result[chatId] = ownerId
	}

	return result, nil
}
