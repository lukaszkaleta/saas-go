package pguser

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	unversalPg "github.com/lukaszkaleta/saas-go/universal/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUser struct {
	Db *pg.PgDb
	Id int64
}

func (pgUser *PgUser) Address() universal.Address {
	return &unversalPg.PgAddress{pgUser.Db, pgUser.tableEntity()}
}

func (pgUser *PgUser) Model() *user.UserModel {
	return &user.UserModel{}
}

func (pgUser *PgUser) Person() universal.Person {
	return &unversalPg.PgPerson{pgUser.Db, pgUser.tableEntity()}
}

func (pgUser *PgUser) Ratings() universal.Ratings {
	return unversalPg.NewPgRatings(pgUser.Db, pgUser.tableEntity())
}

func (pgUser *PgUser) Settings() user.UserSettings {
	return NewPgUserSettings(pgUser.Db, pgUser.Id)
}

func (pgUser *PgUser) Account() user.Account {
	return NewPgAccount(pgUser.Db, pgUser.Id)
}

func (pgUser *PgUser) Archive() error {
	return nil
}

func (pgUser *PgUser) tableEntity() pg.TableEntity {
	return pgUser.Db.TableEntity("users", pgUser.Id)
}
