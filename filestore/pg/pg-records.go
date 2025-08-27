package pgfilestoe

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
)

type PgRecords struct {
	Db *pg.PgDb
}

func (records *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	return nil, nil
}
