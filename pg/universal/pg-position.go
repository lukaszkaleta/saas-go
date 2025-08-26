package universal

import (
	"context"
	"fmt"
	"naborly/internal/api/common"postgres2 "naborly/internal/postgres"

)

type PgPosition struct {
	Db          *pg.PgDb
	TableEntity TableEntity
}

func (pos *PgPosition) Update(model *universal.PositionModel) error {
	query := fmt.Sprintf("update %s set position_latitude = $1, position_longitude = $2 where id = $3", pos.TableEntity.Name)
	_, err := pos.Db.Pool.Exec(context.Background(), query, model.Lat, model.Lon, pos.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pos *PgPosition) Model() *universal.PositionModel {
	return &universal.PositionModel{}
}
