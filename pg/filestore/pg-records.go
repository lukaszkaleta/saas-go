package postgres

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/pg"
)

type PgRecords struct {
	Db *pg.PgDb
}

func (records *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	return nil, nil
}
