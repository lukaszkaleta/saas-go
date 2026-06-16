package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgAction struct {
	db          *pg.PgDb
	tableEntity pg.TableEntity
	name        string
}

func (p PgAction) Model(ctx context.Context) *universal.ActionModel {
	sql := "select action_" + p.name + "_by_id, action_" + p.name + "_at from " + p.tableEntity.Name + " where id = @id"
	rows, err := p.db.Pool.Query(ctx, sql, pgx.NamedArgs{"id": p.tableEntity.Id})
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		var byId *int64
		var madeAt *time.Time
		err := rows.Scan(&byId, &madeAt)
		if err != nil {
			return nil
		}
		return &universal.ActionModel{
			ById:   byId,
			MadeAt: madeAt,
			Name:   p.name,
		}
	}

	return nil
}
