package pg

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func pgContentSetupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "universal-test")
	NewUniversalSchema(db).Create()
	return func(tb testing.TB) {
		NewUniversalSchema(db).Drop()
	}, db
}

func TestPgContents_Add(t *testing.T) {
	teardownSuite, db := pgContentSetupTest(t)
	defer teardownSuite(t)

	contents := PgContents{Db: db}
	value := "{json}"
	content, err := contents.Add(&universal.ContentModel{Value: value, Name: universal.SluggedName("content")})
	if err != nil {
		t.Fatal(err)
	}
	if value != content.Model().Value {
		t.Fatal("Cantent does not match")
	}
}
