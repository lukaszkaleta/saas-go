package pg

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

func ScanMany2many(db *pg.PgDb, query string, ids []*int64) (map[int64]int64, error) {
	rows, err := db.Pool.Query(context.Background(), query, ids)
	if err != nil {
		return nil, err
	}
	var idsMap map[int64]int64
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
