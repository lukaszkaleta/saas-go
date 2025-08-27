package postgres

import (
	"context"
	"naborly/internal/api/user"
)

type PgUserSearch struct {
	Db *PgDb
}

func NewPgUserSearch(db *PgDb) user.UserSearch {
	return &PgUserSearch{Db: db}
}

func (s *PgUserSearch) ByPhone(phone string) (user.User, error) {
	rows, err := s.Db.Pool.Query(context.Background(), "select * from user where person_numner = $1", phone)
	if err != nil {
		return nil, err
	}
	userModel := new(user.UserModel)
	id := 0
	err = UserRowScan(rows, userModel, &id)
	pgUser := &PgUser{Db: s.Db, Id: id}
	if err != nil {
		return pgUser, err
	}
	return user.NewSolidUser(userModel, pgUser, id), nil
}
