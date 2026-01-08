package pguser

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUserSearch struct {
	Db *pg.PgDb
}

func NewPgUserSearch(db *pg.PgDb) user.UserSearch {
	return &PgUserSearch{Db: db}
}

func (s *PgUserSearch) ByPhone(ctx context.Context, phone string) (user.User, error) {
	query := "select * from users where person_phone = @phone"
	rows, err := s.Db.Pool.Query(ctx, query, pgx.NamedArgs{"phone": phone})
	if err != nil {
		return nil, err
	}
	userModels, err := pgx.CollectRows(rows, MapUser)
	if err != nil {
		return nil, err
	}
	if userModels == nil || len(userModels) == 0 {
		return nil, nil
	}
	userModel := userModels[0]
	pgUser := PgUser{Db: s.Db, Id: userModel.Id}
	return user.NewSolidUser(userModel, pgUser), nil
}
