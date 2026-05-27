package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
)

type PgServiceCharges struct {
	db *pg.PgDb
}

func NewPgServiceCharges(db *pg.PgDb) job.ServiceCharges {
	return &PgServiceCharges{db: db}
}

func (p *PgServiceCharges) Active(ctx context.Context) (job.ServiceCharge, error) {
	query := ServiceChargeSelect() + " WHERE active = true"
	rows, err := p.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapServiceCharge(p.db))
}

func (p *PgServiceCharges) Add(ctx context.Context, model *job.ServiceChargeModel) (job.ServiceCharge, error) {
	var id int64
	query := "INSERT INTO service_charge (worker_mode, worker_value, owner_mode, owner_value, active) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := p.db.Pool.QueryRow(ctx, query, model.Worker.Mode, model.Worker.Value, model.Owner.Mode, model.Owner.Value, model.Active).Scan(&id)
	if err != nil {
		return nil, err
	}
	model.Id = id
	return job.NewSolidServiceCharge(model, &PgServiceCharge{db: p.db, Id: id}), nil
}

func (p *PgServiceCharges) All(ctx context.Context) ([]job.ServiceCharge, error) {
	query := ServiceChargeSelect() + " ORDER BY id DESC"
	rows, err := p.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapServiceCharge(p.db))
}
