package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgAddress struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (addr *PgAddress) Update(ctx context.Context, model *universal.AddressModel) error {
	query := fmt.Sprintf("update %s set address_line_1 = $1, address_line_2 = $2, address_city = $3, address_postal_code = $4, address_district = $5 where id = $6", addr.TableEntity.Name)
	_, err := addr.Db.Pool.Exec(ctx, query, model.Line1, model.Line2, model.City, model.PostalCode, model.District, addr.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (addr *PgAddress) Model() *universal.AddressModel {
	return &universal.AddressModel{}
}
