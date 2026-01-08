package pguser

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func pgUsersSetupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "user_test")
	pgfilestore.NewFilestoreSchema(db).Create()
	NewUserSchema(db).Create()

	return func(tb testing.TB) {
		NewUserSchema(db).Drop()
		pgfilestore.NewFilestoreSchema(db).Drop()
	}, db
}

func TestPgUsers_Add(t *testing.T) {
	teardownSuite, db := pgUsersSetupTest(t)
	defer teardownSuite(t)

	users := PgUsers{Db: db}
	personModel := &universal.PersonModel{Phone: "01234"}
	user1, err := users.Add(t.Context(), personModel)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := users.Add(t.Context(), personModel)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Model(t.Context()).Id != user2.Model(t.Context()).Id {
		t.Fatal("Expected same users")
	}
}

func TestPgUsers_ById(t *testing.T) {
	teardownSuite, db := pgUsersSetupTest(t)
	defer teardownSuite(t)

	users := PgUsers{Db: db}
	personModel := &universal.PersonModel{Phone: "01234"}
	user, _ := users.Add(t.Context(), personModel)
	user, err := users.ById(t.Context(), user.Model(t.Context()).Id)
	if err != nil {
		t.Fatal(err)
	}
	if (user.Model(t.Context()).Person.Phone) != personModel.Phone {
		t.Fatal("Expected user with phone")
	}
}
