package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgPerson struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (pgPerson *PgPerson) Ratings() universal.Ratings {
	return NewPgRatings(pgPerson.Db, pgPerson.TableEntity)
}

func (pgPerson *PgPerson) Update(ctx context.Context, model *universal.PersonModel) error {
	query := fmt.Sprintf("update %s set person_first_name = $1, person_last_name = $2, person_email = $3, person_phone = $4 where id = $5", pgPerson.TableEntity.Name)
	_, err := pgPerson.Db.Pool.Exec(ctx, query, model.FirstName, model.LastName, model.Email, model.Phone, pgPerson.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pgPerson *PgPerson) Model(_ context.Context) *universal.PersonModel {
	return &universal.PersonModel{}
}

func (pgPerson *PgPerson) ID() int64 {
	return pgPerson.TableEntity.Id
}

func (pgPerson *PgPerson) Avatar(_ context.Context) universal.Description {
	return NewPgDescription(
		pgPerson.Db,
		pgPerson.TableEntity,
		"avatar_value",
		"avatar_url")
}
