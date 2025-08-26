package universal

import (
	"context"
	"fmt"
	"naborly/internal/api/common"postgres2 "naborly/internal/postgres"

)

type PgPerson struct {
	Db          *pg.PgDb
	TableEntity TableEntity
}

func (pg *PgPerson) Update(model *universal.PersonModel) error {
	query := fmt.Sprintf("update %s set person_first_name = $1, person_last_name = $2, person_email = $3, person_phone = $4 where id = $5", pg.TableEntity.Name)
	_, err := pg.Db.Pool.Exec(context.Background(), query, model.FirstName, model.LastName, model.Email, model.Phone, pg.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgPerson) Model() *universal.PersonModel {
	return &universal.PersonModel{}
}
