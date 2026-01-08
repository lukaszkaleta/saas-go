package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type GlobalJobs interface {
	NearBy(ctx context.Context, position *universal.RadarModel) ([]Job, error)
	ById(ctx context.Context, id int64) (Job, error)
}
