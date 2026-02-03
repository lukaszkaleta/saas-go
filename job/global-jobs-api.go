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

func InvolvedUserIds(ctx context.Context, list []Job) []*int64 {
	idsMap := map[*int64]bool{}
	for _, jobI := range list {
		model := jobI.Model()
		id := model.Actions.Created().ById
		idsMap[id] = true
	}
	ids := make([]*int64, 0, len(idsMap))
	for id := range idsMap {
		ids = append(ids, id)
	}
	return ids
}
