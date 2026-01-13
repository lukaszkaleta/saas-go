package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessage struct {
	Db      *pg.PgDb
	Id      int64
	OwnerId int64
}

func (m *PgMessage) Model(ctx context.Context) *messages.MessageModel {
	query := "select * from message where id=@id"
	rows, err := m.Db.Pool.Query(ctx, query, pgx.NamedArgs{"id": m.Id})
	if err != nil {
		return nil
	}
	model, err := pgx.CollectOneRow(rows, MapMessage)
	if err != nil {
		return nil
	}
	return model
}

func (m *PgMessage) ID() int64 {
	return m.Id
}

func MapMessage(row pgx.CollectableRow) (*messages.MessageModel, error) {
	model := messages.EmptyModel()
	createdActionModel := universal.EmptyCreatedActionModel()
	err := row.Scan(
		&model.Id,
		&model.OwnerId,
		&model.Value,
		&createdActionModel.ById,
		&createdActionModel.MadeAt,
	)
	model.Actions.List[createdActionModel.Name] = createdActionModel
	if err != nil {
		return nil, err
	}
	return model, nil
}
