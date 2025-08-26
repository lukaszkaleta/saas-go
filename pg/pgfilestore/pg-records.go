package pgfilestoe

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/pg/database"
)

type PgRecords struct {
	Db *database.PgDb
}

func (records *PgRecords) Add(ctx context.Context, model *filestore.RecordModel) (filestore.Record, error) {
	return nil, nil
}
