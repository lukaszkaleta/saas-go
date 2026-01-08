package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgPosition struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (pos *PgPosition) Update(ctx context.Context, model *universal.PositionModel) error {
	query := fmt.Sprintf("update %s set position_latitude = $1, position_longitude = $2 where id = $3", pos.TableEntity.Name)
	_, err := pos.Db.Pool.Exec(ctx, query, model.Lat, model.Lon, pos.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pos *PgPosition) Model(ctx context.Context) *universal.PositionModel {
	return &universal.PositionModel{}
}
