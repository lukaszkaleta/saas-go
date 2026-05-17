package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
)

type PgCampaigns struct {
	db *pg.PgDb
}

func NewPgCampaigns(db *pg.PgDb) job.Campaigns {
	return &PgCampaigns{db: db}
}

func (p *PgCampaigns) ById(ctx context.Context, id int64) (job.Campaign, error) {
	query := "SELECT id, description_value, description_image_url, tags FROM campaign WHERE id = $1"
	rows, err := p.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapCampaign(p.db))
}

func (p *PgCampaigns) Add(ctx context.Context, model *job.CampaignModel) (job.Campaign, error) {
	var id int64
	query := "INSERT INTO campaign (description_value, description_image_url, tags) VALUES ($1, $2, $3) RETURNING id"
	err := p.db.Pool.QueryRow(ctx, query, model.Description.Value, model.Description.ImageUrl, model.Tags).Scan(&id)
	if err != nil {
		return nil, err
	}

	pgCampaign := &PgCampaign{
		db: p.db,
		Id: id,
	}

	return job.NewSolidCampaign(&job.CampaignModel{
		Id:          id,
		Description: model.Description,
		Tags:        model.Tags,
	}, pgCampaign), nil
}

func (p *PgCampaigns) List(ctx context.Context) ([]job.Campaign, error) {
	query := "SELECT id, description_value, description_image_url, tags FROM campaign"
	rows, err := p.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapCampaign(p.db))
}

func MapCampaign(db *pg.PgDb) pgx.RowToFunc[job.Campaign] {
	return func(row pgx.CollectableRow) (job.Campaign, error) {
		m, err := MapCampaignModel(row)
		if err != nil {
			return nil, err
		}
		return job.NewSolidCampaign(m, &PgCampaign{db: db, Id: m.Id}), nil
	}
}
