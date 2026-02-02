package pg

import (
	"context"
	"time"

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
	model, err := pgx.CollectOneRow(rows, MapMessageModel)
	if err != nil {
		return nil
	}
	return model
}

func (m *PgMessage) ID() int64 {
	return m.Id
}

func MapMessageModel(row pgx.CollectableRow) (*messages.MessageModel, error) {
	model := messages.EmptyModel()

	actionCreatedModel := universal.EmptyCreatedActionModel()
	actionReadModel := universal.EmptyActionModel("read")

	err := row.Scan(
		&model.Id,
		&model.OwnerId,
		&model.RecipientId,
		&model.Value,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
		&actionReadModel.ById,
		&actionReadModel.MadeAt,
	)
	if err != nil {
		return nil, err
	}
	// created
	model.Actions.List[actionCreatedModel.Name] = actionCreatedModel
	// read
	if actionReadModel.ById == nil {
		noOne := int64(0)
		actionReadModel.ById = &noOne
		actionReadModel.MadeAt = &time.Time{}
	}
	model.Actions.List[actionReadModel.Name] = actionReadModel

	return model, nil
}

func MapMessage(db *pg.PgDb, row pgx.CollectableRow) (messages.Message, error) {
	model, err := MapMessageModel(row)
	if err != nil {
		return nil, err
	}
	pgMessage := &PgMessage{Db: db, Id: model.Id, OwnerId: model.OwnerId}
	return messages.NewSolidMessage(model, pgMessage, model.Id), nil
}
