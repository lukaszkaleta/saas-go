package pguser

import (
	"context"

	"github.com/jackc/pgx/v5"
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
	query := "update users set account_token = $1, firebase_token = $2 where id = $3"
	_, err := pg.Db.Pool.Exec(ctx, query, model.Token, model.FirebaseToken, pg.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgAccount) UpdatePushNotificationToken(ctx context.Context, token string) error {
	query := "update users set firebase_token = @token where id = @id"
	_, err := pg.Db.Pool.Exec(ctx, query, pgx.NamedArgs{"token": &token, "id": pg.Id})
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgAccount) Model(ctx context.Context) (*user.AccountModel, error) {
	query := "select account_token, firebase_token from users where id = $1"
	row := pg.Db.Pool.QueryRow(ctx, query, pg.Id)
	model := &user.AccountModel{}
	err := row.Scan(&model.Token, &model.FirebaseToken)
	if err != nil {
		return nil, err
	}
	return model, nil
}
