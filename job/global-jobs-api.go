package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type GlobalJobs interface {
	universal.FullText[JobSearchOutput]
	Search(ctx context.Context, input *JobSearchInput) ([]*JobSearchOutput, error)
	NearBy(ctx context.Context, position *universal.RadarModel) ([]*JobSearchOutput, error)
	ById(ctx context.Context, id int64) (Job, error)
	ByIds(ctx context.Context, ids []int64) ([]Job, error)
}
