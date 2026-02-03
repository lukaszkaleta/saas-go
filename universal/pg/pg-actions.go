package pg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func NewPgActions(db *pg.PgDb, tableEntity pg.TableEntity) universal.Actions {
	return &PgActions{db: db, tableEntity: tableEntity}
}

type PgActions struct {
	db          *pg.PgDb
	tableEntity pg.TableEntity
}

func (p PgActions) WithName(name string) universal.Action {
	return nil
}

func (p PgActions) Created() universal.Action {
	return p.WithName("created")
}

func (p PgActions) List() map[string]*universal.Action {
	//TODO implement me
	panic("implement me")
}

func (p PgActions) Model(ctx context.Context) (*universal.ActionsModel, error) {
	sql := fmt.Sprintf("select * from %s where id = @id", p.tableEntity.Name)
	rows, err := p.db.Pool.Query(ctx, sql, pgx.NamedArgs{"id": p.tableEntity.Id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actionsMap := make(map[string]*universal.ActionModel)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		fieldDescriptions := rows.FieldDescriptions()

		for i, fd := range fieldDescriptions {
			name := fd.Name
			if strings.HasPrefix(name, "action_") && strings.HasSuffix(name, "_by_id") {
				actionName := strings.TrimSuffix(
					strings.TrimPrefix(name, "action_"),
					"_by_id",
				)
				if actionsMap[actionName] == nil {
					actionsMap[actionName] = universal.EmptyActionModel(actionName)
				}
				id, ok := values[i].(int64)
				if ok {
					actionsMap[actionName].ById = &id
				}
			}
			if strings.HasPrefix(name, "action_") && strings.HasSuffix(name, "_at") {
				actionName := strings.TrimSuffix(
					strings.TrimPrefix(name, "action_"),
					"_at",
				)
				if actionsMap[actionName] == nil {
					actionsMap[actionName] = universal.EmptyActionModel(actionName)
				}
				t, ok := values[i].(time.Time)
				if ok {
					actionsMap[actionName].MadeAt = &t
				}
			}
		}
	}

	return &universal.ActionsModel{List: actionsMap}, err
}
