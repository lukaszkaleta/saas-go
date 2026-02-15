package pg

import (
	"context"
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

const USER_ID = 1

func setupMessagesTest(tb testing.TB) (func(tb testing.TB), context.Context, *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "messages_test")
	schema := NewMessagesSchema(db)
	dropFunc := func(tb testing.TB) {
		err := schema.DropTest()
		if err != nil {
			panic(err)
		}
	}

	dropFunc(tb)
	err := schema.CreateTest()

	if err != nil {
		dropFunc(tb)
		tb.Fatal(err)
	}

	for i := 1; i < 3; i++ {
		_, err := db.Pool.Exec(tb.Context(), "insert into users (id) values ($1)", i)
		if err != nil {
			tb.Error(err)
		}
	}

	ctx := user.WithUser(tb.Context(), user.SolidUser{Id: USER_ID})
	return dropFunc, ctx, db
}

func TestPgMessages_Add(t *testing.T) {
	drop, ctx, db := setupMessagesTest(t)
	defer drop(t)

	pgMessages := NewPgMessages(db, pg.RelationEntity{TableName: "message", RelationId: USER_ID})
	value := "test-message"
	newMessage, err := pgMessages.Add(ctx, USER_ID, value)
	if err != nil {
		t.Fatal(err)
	}
	if newMessage == nil {
		t.Fatal("newMessage is nil")
	}
	model, err := newMessage.Model(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if model.OwnerId != USER_ID {
		t.Fatal("Wrong Owner id")
	}
	messageModel, err := newMessage.Model(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if messageModel.Value != value {
		t.Fatal("Wrong value")
	}

	list, err := pgMessages.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatal("Wrong list length")
	}
}
