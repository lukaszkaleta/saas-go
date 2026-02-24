package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Tasks interface {
	Create(ctx context.Context, model *TaskModel) (Task, error)
	ById(ctx context.Context, id int64) (Task, error)
	ByJobId(ctx context.Context, jobId int64) (Task, error)
	Current(ctx context.Context) ([]Task, error)
	InProgress(ctx context.Context) ([]Task, error)
	Completed(ctx context.Context) ([]Task, error)
	WaitingForPayment(ctx context.Context) ([]Task, error)
	Earnings(ctx context.Context) (map[string]*universal.PriceModel, error)
}
