package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgArrayCounter struct {
	Db       *pg.PgDb
	Relation pg.RelationEntity
	GetId    func() int64
}

func NewPgArrayCounter(db *pg.PgDb, relation pg.RelationEntity, getId func() int64) universal.Counter {
	return &PgArrayCounter{Db: db, Relation: relation, GetId: getId}
}

func (c *PgArrayCounter) Increment(ctx context.Context) error {
	id := c.GetId()
	query := fmt.Sprintf("UPDATE %s SET %s = array_append(%s, @item_id) WHERE id = @id AND NOT (@item_id = ANY(%s))", c.Relation.TableName, c.Relation.ColumnName, c.Relation.ColumnName, c.Relation.ColumnName)
	_, err := c.Db.Pool.Exec(ctx, query, pgx.NamedArgs{"item_id": id, "id": c.Relation.RelationId})
	return err
}

func (c *PgArrayCounter) Get(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT cardinality(%s) FROM %s WHERE id = @id", c.Relation.ColumnName, c.Relation.TableName)
	err := c.Db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"id": c.Relation.RelationId}).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
