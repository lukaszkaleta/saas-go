package pguser

import (
	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
	unversalPg "github.com/lukaszkaleta/saas-go/universal/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUser struct {
	Db *pg.PgDb
	Id int64
}

func (pgUser PgUser) ID() int64 {
	return pgUser.Id
}

func (pgUser PgUser) Address() universal.Address {
	return &unversalPg.PgAddress{pgUser.Db, pgUser.TableEntity()}
}

func (pgUser PgUser) Model() *user.UserModel {
	return &user.UserModel{}
}

func (pgUser PgUser) Person() universal.Person {
	return &unversalPg.PgPerson{pgUser.Db, pgUser.TableEntity()}
}

func (pgUser PgUser) Settings() user.UserSettings {
	return NewPgUserSettings(pgUser.Db, pgUser.Id)
}

func (pgUser PgUser) Account() user.Account {
	return NewPgAccount(pgUser.Db, pgUser.Id)
}

func (pgUser PgUser) Archive() error {
	return nil
}

func (pgUser PgUser) FileSystem(name string) (filestore.FileSystem, error) {
	return nil, nil
}

func (pgUser PgUser) TableEntity() pg.TableEntity {
	return pgUser.Db.TableEntity("users", pgUser.Id)
}

func MapUser(row pgx.CollectableRow) (*user.UserModel, error) {
	userModel := user.EmptyUserModel()
	err := row.Scan(
		&userModel.Id,
		&userModel.Account.Token,
		&userModel.Person.FirstName,
		&userModel.Person.LastName,
		&userModel.Person.Email,
		&userModel.Person.Phone,
		&userModel.Address.Line1,
		&userModel.Address.Line2,
		&userModel.Address.City,
		&userModel.Address.PostalCode,
		&userModel.Address.District,
		&userModel.Settings.Avatar.Value,
		&userModel.Settings.Avatar.ImageUrl,
		&userModel.Settings.Radar.Perimeter,
		&userModel.Settings.Radar.Position.Lat,
		&userModel.Settings.Radar.Position.Lon,
	)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}
