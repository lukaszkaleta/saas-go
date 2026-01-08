package pguser

import (
	"context"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgAccount struct {
	Db *pg.PgDb
	Id int64
}

func NewPgAccount(Db *pg.PgDb, id int64) user.Account {
	return &PgAccount{Db: Db, Id: id}
}

func (pg *PgAccount) Update(ctx context.Context, model *user.AccountModel) error {
	query := "update users set account_token = $1 where id = $1"
	_, err := pg.Db.Pool.Exec(ctx, query, model.Token, pg.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgAccount) Model(ctx context.Context) *user.AccountModel {
	return &user.AccountModel{}
}
