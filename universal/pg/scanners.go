package pg

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

func ScanMany2many(db *pg.PgDb, ctx context.Context, query string, ids []int64) (map[int64]int64, error) {
	rows, err := db.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	idsMap := make(map[int64]int64)
	for rows.Next() {
		var idKey int64
		var idValue int64
		err = rows.Scan(&idKey, &idValue)
		if err != nil {
			return nil, err
		}
		idsMap[idKey] = idValue
	}
	rows.Close()
	return idsMap, nil
}
