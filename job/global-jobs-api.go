package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type GlobalJobs interface {
	universal.FullText[Job]
	Search(ctx context.Context, input JobSearchInput) ([]Job, error)
	NearBy(ctx context.Context, position *universal.RadarModel) ([]Job, error)
	ById(ctx context.Context, id int64) (Job, error)
}
