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
	userId := int64(0)
	query := "INSERT INTO users(person_first_name, person_last_name, person_email, person_phone) VALUES( $1, $2, $3, $4 ) returning id"
	row := pgUsers.Db.Pool.QueryRow(context.Background(), query, model.FirstName, model.LastName, model.Email, model.Phone)
	row.Scan(&userId)
	pgUser := &PgUser{
		Db: pgUsers.Db,
		Id: userId,
	}
	return user.NewSolidUser(
		&user.UserModel{Id: userId, Person: model, Address: &universal.AddressModel{}},
		pgUser,
		userId,
	), nil
}

func (pgUsers *PgUsers) ById(id int64) (user.User, error) {
	personRow := new(universal.PersonModel)
	accountRow := new(user.AccountModel)
	addressRow := new(universal.AddressModel)
	query := "select * from users where id = $1"
	row := pgUsers.Db.Pool.QueryRow(context.Background(), query, id)
	err := row.Scan(
		&id,
		&accountRow.Token,
		&personRow.FirstName,
		&personRow.FirstName,
		&personRow.LastName,
		&personRow.Email,
		&personRow.Phone,
		&addressRow.Line1,
		&addressRow.Line2,
		&addressRow.City,
		&addressRow.PostalCode,
		&addressRow.District,
	)
	pgUser := &PgUser{Db: pgUsers.Db, Id: id}
	if err != nil {
		return pgUser, err
	}

	return user.NewSolidUser(
		&user.UserModel{Id: id, Person: personRow, Address: addressRow},
		pgUser,
		id,
	), nil
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
