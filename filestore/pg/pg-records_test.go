package pgfilestore

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

func TestPgRecords_Urls(t *testing.T) {
	teardownSuite, db := setupPgFilestoreTest(t)
	defer teardownSuite(t)

	relationEntity := pg.RelationEntity{
		TableName:  "test_filesystem",
		RelationId: 10,
		ColumnName: "test_id",
	}
	pgRecords := NewPgRecords(db, NewPgFileSystem(db, relationEntity))
	records := pgRecords
	_, err := records.AddFromUrl(t.Context(), "url1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = records.AddFromUrl(t.Context(), "url2")
	if err != nil {
		t.Fatal(err)
	}
	urls, err := pgRecords.Urls(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 2 {
		t.Errorf("got %d urls, want 2", len(urls))
	}
}
