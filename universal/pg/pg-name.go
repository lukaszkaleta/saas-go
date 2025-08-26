package postgres

import (
	"context"
	"fmt"
	"naborly/internal/api/common"postgres2 "naborly/internal/postgres"

)

type PgName struct {
	Db          *postgres2.PgDb
	TableEntity TableEntity
}

func (p *PgName) Update(model *common.NameModel) error {
	query := fmt.Sprintf("update %s set name_value = $1, name_slug = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Value, model.Slug, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgName) Model() *common.NameModel {
	return &common.NameModel{}
}
