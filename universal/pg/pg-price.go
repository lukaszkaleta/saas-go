package postgres

import (
	"context"
	"fmt"
	"naborly/internal/api/common"
	postgres2 "naborly/internal/postgres"
)

type PgPrice struct {
	Db          *postgres2.PgDb
	TableEntity TableEntity
}

func (p *PgPrice) Update(model *common.PriceModel) error {
	query := fmt.Sprintf("update %s set price_value = $1, price_currency = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Value, model.Currency, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgPrice) Model() *common.PriceModel {
	return &common.PriceModel{}
}
