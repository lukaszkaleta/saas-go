package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Campaign interface {
	universal.Idable
	Model(ctx context.Context) (*CampaignModel, error)
	Description() universal.Description
}

type CampaignModel struct {
	Id          int64                       `json:"id"`
	Description *universal.DescriptionModel `json:"description"`
	Tags        []string                    `json:"tags"`
}

func (m CampaignModel) ID() int64 {
	return m.Id
}

type SolidCampaign struct {
	model    *CampaignModel
	campaign Campaign
}

func NewSolidCampaign(model *CampaignModel, campaign Campaign) Campaign {
	return &SolidCampaign{model: model, campaign: campaign}
}

func (s *SolidCampaign) ID() int64 {
	return s.model.Id
}

func (s *SolidCampaign) Model(ctx context.Context) (*CampaignModel, error) {
	return s.model, nil
}

func (s *SolidCampaign) Description() universal.Description {
	return s.campaign.Description()
}
