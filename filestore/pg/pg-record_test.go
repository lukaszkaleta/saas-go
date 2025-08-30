package pgfilestoe

import (
	"log"
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test suite")
	return func(tb testing.TB) {
		log.Println("teardown test suite")
	}
}

func setupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.NewPg()
	schema := NewFilestoreSchema(db)
	schema.Create()

	return func(tb testing.TB) {
		schema.Drop()
	}, db
}

func TestPgRecord_Update(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	records := PgRecords{Db: db}
	record, err := records.Add(t.Context(), filestore.EmptyRecordModel())
	if err != nil {
		t.Fatal(err)
	}
	err = record.Update(&filestore.RecordModel{Name: &universal.NameModel{Value: "koza"}, Description: &universal.DescriptionModel{}})
	if err != nil {
		t.Fatal(err)
	}
}
