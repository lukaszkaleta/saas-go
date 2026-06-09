package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
)

type PgSimpleCounter struct {
	Db       *pg.PgDb
	Relation pg.RelationEntity
}

func (c *PgSimpleCounter) Increment(ctx context.Context) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s + 1 WHERE id = @id", c.Relation.TableName, c.Relation.ColumnName, c.Relation.ColumnName)
	_, err := c.Db.Pool.Exec(ctx, query, pgx.NamedArgs{"id": c.Relation.RelationId})
	return err
}

func (c *PgSimpleCounter) Get(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = @id", c.Relation.ColumnName, c.Relation.TableName)
	err := c.Db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"id": c.Relation.RelationId}).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
