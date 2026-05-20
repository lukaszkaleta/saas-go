package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
)

type PgCampaignJobs struct {
	db *pg.PgDb
}

func NewPgCampaignJobs(db *pg.PgDb) job.CampaignJobs {
	return &PgCampaignJobs{db: db}
}

func (p *PgCampaignJobs) TopActive(ctx context.Context) ([]job.CampaignJob, error) {
	// 1. Read campaigns from campaign table
	campaignQuery := "SELECT id, description_value, description_image_url, tags FROM campaign"
	campaignRows, err := p.db.Pool.Query(ctx, campaignQuery)
	if err != nil {
		return nil, err
	}
	campaignModels, err := pgx.CollectRows(campaignRows, MapCampaignModel)
	if err != nil {
		return nil, err
	}

	result := make([]job.CampaignJob, 0, len(campaignModels))

	// 2. For each campaign, read jobs by tags (at least one tag should match)
	for _, cm := range campaignModels {
		jobQuery := JobSelect() + " WHERE tags && @tags AND " + WhereStatusIsPublic()
		jobRows, err := p.db.Pool.Query(ctx, jobQuery, pgx.NamedArgs{"tags": cm.Tags})
		if err != nil {
			return nil, err
		}
		jobs, err := pgx.CollectRows(jobRows, MapJob(p.db))
		if err != nil {
			return nil, err
		}

		result = append(result, job.CampaignJob{
			Id:   cm.Id,
			Name: cm.Description.Value,
			Jobs: jobs,
		})
	}

	return result, nil
}
