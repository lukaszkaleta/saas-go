package universal

import (
	"context"
	"fmt"
	"saas-go/pg"
	"saas-go/universal"
)

type PgPrice struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (p *PgPrice) Update(model *universal.PriceModel) error {
	query := fmt.Sprintf("update %s set price_value = $1, price_currency = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Value, model.Currency, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgPrice) Model() *universal.PriceModel {
	return &universal.PriceModel{}
}
