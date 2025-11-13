package job

import (
	"context"
	"fmt"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type GlobalJobsLogger struct {
	next GlobalJobs
}

func (g *GlobalJobsLogger) ById(ctx context.Context, id int64) (Job, error) {
	//TODO implement me
	panic("implement me")
}

func NewGlobalJobsLogger(next GlobalJobs) GlobalJobs {
	return &GlobalJobsLogger{
		next: next,
	}
}

func (g *GlobalJobsLogger) NearBy(position *universal.RadarModel) ([]Job, error) {
	defer func(start time.Time) {
		fmt.Printf("Searching for jobs took %v\n", time.Since(start))
	}(time.Now())
	return g.next.NearBy(position)
}
