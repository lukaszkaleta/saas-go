package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/pg/database"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgName struct {
	Db          *database.PgDb
	TableEntity database.TableEntity
}

func (p *PgName) Update(model *universal.NameModel) error {
	query := fmt.Sprintf("update %s set name_value = $1, name_slug = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Value, model.Slug, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgName) Model() *universal.NameModel {
	return &universal.NameModel{}
}
