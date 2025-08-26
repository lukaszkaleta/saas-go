package pg_universal

import (
	"context"
	"fmt"
	"github.com/lukaszkaleta/saas-go/pg/database"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgDescription struct {
	Db          *database.PgDb
	TableEntity database.TableEntity
}

func (p *PgDescription) Update(model *universal.DescriptionModel) error {
	query := fmt.Sprintf("update %s set description_value = $1, description_image_url = $2 where id = $3", p.TableEntity.Name)
	_, err := p.Db.Pool.Exec(context.Background(), query, model.Value, model.ImageUrl, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgDescription) Model() *universal.DescriptionModel {
	return &universal.DescriptionModel{}
}
