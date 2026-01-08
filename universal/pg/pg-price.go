package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgPrice struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (p *PgPrice) Update(ctx context.Context, model *universal.PriceModel) error {
	query := fmt.Sprintf("update %s set price_value = $1, price_currency = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(ctx, query, model.Value, model.Currency, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgPrice) Model(ctx context.Context) *universal.PriceModel {
	return &universal.PriceModel{}
}
