package pguser

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUsers struct {
	Db *pg.PgDb
}

func NewPgUsers(s *pg.PgDb) user.Users {
	return &PgUsers{Db: s}
}

func (pgUsers *PgUsers) Search() user.UserSearch {
	return NewPgUserSearch(pgUsers.Db)
}

func (pgUsers *PgUsers) Add(model *universal.PersonModel) (user.User, error) {
	userWithPhone, err := pgUsers.Search().ByPhone(model.Phone)
	if err != nil {
		return nil, err
	}
	if userWithPhone != nil {
		return userWithPhone, nil
	}

	userId := int64(0)
	query := "INSERT INTO users(person_first_name, person_last_name, person_email, person_phone) VALUES( $1, $2, $3, $4 ) returning id"
	row := pgUsers.Db.Pool.QueryRow(context.Background(), query, model.FirstName, model.LastName, model.Email, model.Phone)
	row.Scan(&userId)
	pgUser := &PgUser{
		Db: pgUsers.Db,
		Id: userId,
	}
	userModel := user.EmptyUserModel()
	userModel.Person = model
	return user.NewSolidUser(
		userModel,
		pgUser,
		userId,
	), nil
}

func (pgUsers *PgUsers) ById(id int64) (user.User, error) {
	sql := "select * from users where id = @id"
	rows, err := pgUsers.Db.Pool.Query(context.Background(), sql, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	pgUser := PgUser{
		Db: pgUsers.Db,
		Id: id,
	}
	userModel, err := pgx.CollectOneRow(rows, MapUser)
	return user.NewSolidUser(userModel, pgUser, id), err
}

func (pgUsers *PgUsers) ListAll() ([]user.User, error) {
	query := "select * from users"
	rows, err := pgUsers.Db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	users := []user.User{}
	for rows.Next() {
		userModel := user.NewUserModel()
		id := int64(0)
		err := UserRowScan(rows, userModel, &id)
		if err != nil {
			return nil, err
		}
		pgUser := PgUser{Db: pgUsers.Db, Id: id}
		solidUser := user.NewSolidUser(userModel, &pgUser, id)
		users = append(users, solidUser)
	}
	return users, nil
}

func (pgUsers *PgUsers) EstablishAccount(model *user.UserModel) (user.User, error) {
	user, err := pgUsers.Search().ByPhone(model.Person.Phone)
	if err != nil {
		return user, err
	}

	if user != nil {
		user.Person().Update(model.Person)
	} else {
		user, err = pgUsers.Add(model.Person)
		if err != nil {
			return user, err
		}
	}

	user.Account().Update(model.Account)
	user.Address().Update(model.Address)

	return user, nil
}

func UserRowScan(row pgx.Rows, userRow *user.UserModel, id *int64) error {
	return row.Scan(
		&id,
		&userRow.Account.Token,
		&userRow.Person.FirstName,
		&userRow.Person.LastName,
		&userRow.Person.Email,
		&userRow.Person.Phone,
		&userRow.Address.Line1,
		&userRow.Address.Line2,
		&userRow.Address.City,
		&userRow.Address.PostalCode,
		&userRow.Address.District,
		&userRow.Settings.Radar.Perimeter,
		&userRow.Settings.Radar.Position.Lon,
		&userRow.Settings.Radar.Position.Lat,
	)
}
