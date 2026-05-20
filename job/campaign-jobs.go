package job

import (
	"context"
)

type CampaignJobs interface {
	TopActive(ctx context.Context) ([]CampaignJob, error)
}
