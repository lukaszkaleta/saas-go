package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

type PgCounter struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
	ColumnName  string
}

func (c *PgCounter) Increment(ctx context.Context) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s + 1 WHERE id = $1", c.TableEntity.Name, c.ColumnName, c.ColumnName)
	_, err := c.Db.Pool.Exec(ctx, query, c.TableEntity.Id)
	return err
}

func (c *PgCounter) Get(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", c.ColumnName, c.TableEntity.Name)
	err := c.Db.Pool.QueryRow(ctx, query, c.TableEntity.Id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
