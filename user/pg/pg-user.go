package pguser

import (
	"context"
	"strings"

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

func (pgUser PgUser) Model(ctx context.Context) *user.UserModel {
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

func (pgUser PgUser) Rated() universal.Rated {
	return unversalPg.NewPgRated(pgUser.Db, pgUser.Id, "job_rating")
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

func MapUser(db *pg.PgDb) pgx.RowToFunc[user.User] {
	return func(row pgx.CollectableRow) (user.User, error) {
		model, err := MapUserModel(row)
		if err != nil {
			return nil, err
		}
		pgUser := PgUser{Db: db, Id: model.Id}
		return user.NewSolidUser(model, pgUser), nil
	}
}

func MapUserModel(row pgx.CollectableRow) (*user.UserModel, error) {
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
		&userModel.Person.Avatar.Value,
		&userModel.Person.Avatar.ImageUrl,
		&userModel.Settings.Radar.Perimeter,
		&userModel.Settings.Radar.Position.Lat,
		&userModel.Settings.Radar.Position.Lon,
		&userModel.Person.AverageRating,
	)
	if err != nil {
		return nil, err
	}
	userModel.Person.Id = userModel.Id
	return userModel, nil
}

func UserColumns() []string {
	return []string{
		"id",
		"account_token",
		"person_first_name",
		"person_last_name",
		"person_email",
		"person_phone",
		"address_line_1",
		"address_line_2",
		"address_city",
		"address_postal_code",
		"address_district",
		"avatar_description_value",
		"avatar_description_image_url",
		"settings_radar_perimeter",
		"settings_radar_position_latitude",
		"settings_radar_position_longitude",
		"ratings_average",
	}
}

func UserColumnString() string {
	return strings.Join(UserColumns(), ",")
}

func UserColumnsSelect() string {
	return "select " + UserColumnString()
}

func UserSelect() string {
	return UserColumnsSelect() + " from users "
}
