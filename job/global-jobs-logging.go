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

func NewGlobalJobsLogger(next GlobalJobs) GlobalJobs {
	return &GlobalJobsLogger{
		next: next,
	}
}

func (g *GlobalJobsLogger) ById(ctx context.Context, id int64) (Job, error) {
	return g.next.ById(ctx, id)
}

func (g *GlobalJobsLogger) NearBy(ctx context.Context, position *universal.RadarModel) ([]JobSearchOutput, error) {
	defer func(start time.Time) {
		fmt.Printf("Searching for jobs took %v\n", time.Since(start))
	}(time.Now())
	return g.next.NearBy(ctx, position)
}

func (g *GlobalJobsLogger) ByQuery(ctx context.Context, query *string) ([]JobSearchOutput, error) {
	defer func(start time.Time) {
		fmt.Printf("ByQuery for jobs took %v\n", time.Since(start))
	}(time.Now())
	return g.next.ByQuery(ctx, query)
}

func (g *GlobalJobsLogger) Search(ctx context.Context, input JobSearchInput) ([]JobSearchOutput, error) {
	defer func(start time.Time) {
		fmt.Printf("ByQuery for jobs took %v\n", time.Since(start))
	}(time.Now())
	return g.next.Search(ctx, input)
}
