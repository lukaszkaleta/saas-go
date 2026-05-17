package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgCampaign struct {
	db *pg.PgDb
	Id int64
}

func (p *PgCampaign) ID() int64 {
	return p.Id
}

func (p *PgCampaign) Model(ctx context.Context) (*job.CampaignModel, error) {
	query := "SELECT id, description_value, description_image_url, tags FROM campaign WHERE id = $1"
	rows, err := p.db.Pool.Query(ctx, query, p.Id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapCampaignModel)
}

func (p *PgCampaign) Description() universal.Description {
	return pgUniversal.NewPgDescriptionFromTable(p.db, p.tableEntity())
}

func (p *PgCampaign) tableEntity() pg.TableEntity {
	return p.db.TableEntity("campaign", p.Id)
}

func MapCampaignModel(row pgx.CollectableRow) (*job.CampaignModel, error) {
	m := &job.CampaignModel{
		Description: &universal.DescriptionModel{},
	}
	err := row.Scan(&m.Id, &m.Description.Value, &m.Description.ImageUrl, &m.Tags)
	if err != nil {
		return nil, err
	}
	return m, nil
}
