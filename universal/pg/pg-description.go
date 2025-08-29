package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgDescription struct {
	tableEntity pg.TableEntity
	db          *pg.PgDb
	valueColumn string
	urlColumn   string
}

func NewPgDescriptionFromDb(db *pg.PgDb, id int64) *PgDescription {
	return NewPgDescriptionFromTable(db, pg.TableEntity{Name: "description", Id: id})
}

func NewPgDescriptionFromTable(db *pg.PgDb, entity pg.TableEntity) *PgDescription {
	return NewPgDescription(db, entity, "description_value", "description_url")
}

func NewPgDescription(db *pg.PgDb, entity pg.TableEntity, valueColumn string, urlColumn string) *PgDescription {
	return &PgDescription{db: db, tableEntity: entity, valueColumn: valueColumn, urlColumn: urlColumn}
}

func (p *PgDescription) Update(model *universal.DescriptionModel) error {
	query := fmt.Sprintf("update %s set %s = $1, %s = $2 where id = $3", p.tableEntity.Name, p.valueColumn, p.urlColumn)
	_, err := p.db.Pool.Exec(context.Background(), query, model.Value, model.ImageUrl, p.tableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgDescription) Model() *universal.DescriptionModel {
	return &universal.DescriptionModel{}
}

func UseMapDescription(model *universal.DescriptionModel) MapName {
	return func(row pgx.CollectableRow) {
		err := row.Scan(&model.Value)
		if err != nil {
			fmt.Printf("Error scanning description value: %v\n", err)
		}
		err = row.Scan(&model.ImageUrl)
		if err != nil {
			fmt.Printf("Error scanning description image url: %v\n", err)
		}
	}
}

type MapDescription func(row pgx.CollectableRow)
