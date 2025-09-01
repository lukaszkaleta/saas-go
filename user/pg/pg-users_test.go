package pguser

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func setupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("users_test")
	pgfilestore.NewFilestoreSchema(db).Create()
	NewUserSchema(db).Create()

	return func(tb testing.TB) {
		NewUserSchema(db).Drop()
		pgfilestore.NewFilestoreSchema(db).Drop()
	}, db
}

func TestPgUsers_Add(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	users := PgUsers{Db: db}
	personModel := &universal.PersonModel{Phone: "01234"}
	user1, err := users.Add(personModel)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := users.Add(personModel)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Model().Id != user2.Model().Id {
		t.Fatal("Expected same users")
	}
}
