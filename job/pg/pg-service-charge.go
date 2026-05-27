package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
)

type PgServiceCharge struct {
	db *pg.PgDb
	Id int64
}

func (p *PgServiceCharge) ID() int64 {
	return p.Id
}

func (p *PgServiceCharge) Model(ctx context.Context) (*job.ServiceChargeModel, error) {
	query := ServiceChargeSelect() + " WHERE id = $1"
	rows, err := p.db.Pool.Query(ctx, query, p.Id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapServiceChargeModel)
}

func (p *PgServiceCharge) Activate(ctx context.Context) error {
	_, err := p.db.Pool.Exec(ctx, "UPDATE service_charge SET active = (id = $1)", p.Id)
	return err
}

func (p *PgServiceCharge) Deactivate(ctx context.Context) error {
	_, err := p.db.Pool.Exec(ctx, "UPDATE service_charge SET active = false WHERE id = $1", p.Id)
	return err
}

func (p *PgServiceCharge) IsActive(ctx context.Context) (bool, error) {
	var active bool
	err := p.db.Pool.QueryRow(ctx, "SELECT active FROM service_charge WHERE id = $1", p.Id).Scan(&active)
	return active, err
}

func MapServiceCharge(db *pg.PgDb) pgx.RowToFunc[job.ServiceCharge] {
	return func(row pgx.CollectableRow) (job.ServiceCharge, error) {
		model, err := MapServiceChargeModel(row)
		if err != nil {
			return nil, err
		}
		return job.NewSolidServiceCharge(model, &PgServiceCharge{db: db, Id: model.Id}), nil
	}
}

func MapServiceChargeModel(row pgx.CollectableRow) (*job.ServiceChargeModel, error) {
	var m job.ServiceChargeModel
	err := row.Scan(&m.Id, &m.Worker.Mode, &m.Worker.Value, &m.Owner.Mode, &m.Owner.Value, &m.Active)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func ServiceChargeColumns() []string {
	return []string{"id", "worker_mode", "worker_value", "owner_mode", "owner_value", "active"}
}

func ServiceChargeSelect() string {
	return "SELECT id, worker_mode, worker_value, owner_mode, owner_value, active FROM service_charge"
}
