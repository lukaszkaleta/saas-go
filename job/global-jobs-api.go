package job

import "github.com/lukaszkaleta/saas-go/universal"

type GlobalJobs interface {
	NearBy(position *universal.RadarModel) ([]Job, error)
	ById(id int64) (Job, error)
}
