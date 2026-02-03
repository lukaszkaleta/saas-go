package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgMessage struct {
	Db    *pg.PgDb
	Id    int64
	Owner pg.RelationEntity
}

func (m *PgMessage) Acknowledge(ctx context.Context) error {
	currentUserId := universal.CurrentUserId(ctx)
	query := fmt.Sprintf("update %s set action_read_at = now(), action_read_by_id = @userId where id = @id", m.Owner.TableName)
	_, err := m.Db.Pool.Exec(ctx, query, pgx.NamedArgs{"userId": currentUserId, "id": m.Id})
	return err
}

func (m *PgMessage) Model(ctx context.Context) *messages.MessageModel {
	query := fmt.Sprintf(ColumnsSelect()+" from %s where id=@id", m.Owner.TableName)
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
	model.Actions.List[actionCreatedModel.Name] = actionCreatedModel
	model.Actions.List[actionReadModel.Name] = actionReadModel

	return model, nil
}

func MapMessage(db *pg.PgDb, owner pg.RelationEntity) pgx.RowToFunc[messages.Message] {
	return func(row pgx.CollectableRow) (messages.Message, error) {
		model, err := MapMessageModel(row)
		if err != nil {
			return nil, err
		}
		pgMessage := &PgMessage{Db: db, Id: model.Id, Owner: owner}
		return messages.NewSolidMessage(model, pgMessage, model.Id), nil
	}
}

func Columns() []string {
	return []string{
		"id",
		"owner_id",
		"recipient_id",
		"value",
		"action_created_by_id",
		"action_created_at",
		"action_read_by_id",
		"action_read_at",
	}
}

func ColumnString() string {
	return strings.Join(Columns(), ",")
}

func ColumnsSelect() string {
	return "select " + ColumnString()
}

func Select(table string) string {
	return ColumnsSelect() + " from " + table + " "
}
