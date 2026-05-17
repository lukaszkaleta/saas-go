package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Campaigns interface {
	universal.Idables[Campaign]
	Add(ctx context.Context, model *CampaignModel) (Campaign, error)
	List(ctx context.Context) ([]Campaign, error)
}

func CampaignModels(ctx context.Context, list []Campaign) []*CampaignModel {
	models := make([]*CampaignModel, len(list))
	for i, c := range list {
		models[i], _ = c.Model(ctx)
	}
	return models
}
