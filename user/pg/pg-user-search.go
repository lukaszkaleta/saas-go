package pguser

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUserSearch struct {
	Db *pg.PgDb
}

func NewPgUserSearch(db *pg.PgDb) user.UserSearch {
	return &PgUserSearch{Db: db}
}

func (s *PgUserSearch) ByPhone(phone string) (user.User, error) {
	rows, err := s.Db.Pool.Query(context.Background(), "select * from user where person_numner = $1", phone)
	if err != nil {
		return nil, err
	}
	userModel := new(user.UserModel)
	id := int64(0)
	err = UserRowScan(rows, userModel, &id)
	pgUser := &PgUser{Db: s.Db, Id: id}
	if err != nil {
		return pgUser, err
	}
	return user.NewSolidUser(userModel, pgUser, id), nil
}
