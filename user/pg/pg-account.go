package postgres

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

func (pg *PgAccount) Update(model *user.AccountModel) error {
	query := "update users set account_token = $1 where id = $1"
	_, err := pg.Db.Pool.Exec(context.Background(), query, model.Token, pg.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgAccount) Model() *user.AccountModel {
	return &user.AccountModel{}
}
