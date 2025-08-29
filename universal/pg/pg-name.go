package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgName struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
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

func UseMapName(nameModel *universal.NameModel) MapName {
	return func(row pgx.CollectableRow) {
		err := row.Scan(&nameModel.Value)
		if err != nil {
			fmt.Printf("Error scanning name: %v\n", err)
		}
		err = row.Scan(&nameModel.Slug)
		if err != nil {
			fmt.Printf("Error scanning name: %v\n", err)
		}
	}
}

type MapName func(row pgx.CollectableRow)
