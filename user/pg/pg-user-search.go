package pguser

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
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
	return pgx.CollectOneRow(rows, MapUser(s.Db))
}

func (s *PgUserSearch) ModelsByIds(ctx context.Context, ids []*int64) ([]*user.UserModel, error) {
	query := "select * from users where id = any(@ids)"
	rows, err := s.Db.Pool.Query(ctx, query, pgx.NamedArgs{"ids": ids})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapUserModel)
}

func (s *PgUserSearch) PersonModelsByIds(ctx context.Context, ids []*int64) ([]*universal.PersonModel, error) {
	userModesl, err := s.ModelsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	personModels := make([]*universal.PersonModel, len(userModesl))
	for i, userModel := range userModesl {
		personModels[i] = userModel.Person
	}
	return personModels, nil
}
