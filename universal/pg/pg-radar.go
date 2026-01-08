package pg

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgRadar struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (r *PgRadar) Update(ctx context.Context, model *universal.RadarModel) error {
	// update perimeter and position columns on the target table
	query := fmt.Sprintf("update %s set settings_radar_perimeter = $1, settings_radar_position_latitude = $2, settings_radar_position_longitude = $3 where id = $4", r.TableEntity.Name)
	_, err := r.Db.Pool.Exec(ctx, query, model.Perimeter, model.Position.Lat, model.Position.Lon, r.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgRadar) Model(ctx context.Context) *universal.RadarModel {
	return &universal.RadarModel{}
}
