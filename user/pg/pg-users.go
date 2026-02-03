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

func (pgUsers *PgUsers) Add(ctx context.Context, model *universal.PersonModel) (user.User, error) {
	userWithPhone, err := pgUsers.Search().ByPhone(ctx, model.Phone)
	if err != nil {
		return nil, err
	}
	if userWithPhone != nil {
		return userWithPhone, nil
	}

	userId := int64(0)
	query := "INSERT INTO users(person_first_name, person_last_name, person_email, person_phone) VALUES( $1, $2, $3, $4 ) returning id"
	row := pgUsers.Db.Pool.QueryRow(ctx, query, model.FirstName, model.LastName, model.Email, model.Phone)
	row.Scan(&userId)
	pgUser := &PgUser{
		Db: pgUsers.Db,
		Id: userId,
	}
	userModel := user.EmptyUserModel()
	userModel.Person = model
	userModel.Id = userId
	return user.NewSolidUser(
		userModel,
		pgUser,
	), nil
}

func (pgUsers *PgUsers) ById(ctx context.Context, id int64) (user.User, error) {
	sql := "select * from users where id = @id"
	rows, err := pgUsers.Db.Pool.Query(ctx, sql, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapUser(pgUsers.Db))
}

func (pgUsers *PgUsers) ListAll(ctx context.Context) ([]user.User, error) {
	query := "select * from users"
	rows, err := pgUsers.Db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapUser(pgUsers.Db))
}

func (pgUsers *PgUsers) EstablishAccount(ctx context.Context, model *user.UserModel) (user.User, error) {
	userByPhone, err := pgUsers.Search().ByPhone(ctx, model.Person.Phone)
	if err != nil {
		return userByPhone, err
	}

	if userByPhone != nil {
		userByPhone.Person().Update(ctx, model.Person)
	} else {
		userByPhone, err = pgUsers.Add(ctx, model.Person)
		if err != nil {
			return userByPhone, err
		}
	}

	userByPhone.Account().Update(ctx, model.Account)
	userByPhone.Address().Update(ctx, model.Address)

	return userByPhone, nil
}
