package job

import "context"

type Tasks interface {
	Create(ctx context.Context, model *TaskModel) (Task, error)
	ByJobId(ctx context.Context, jobId int64) (Task, error)
	Current(ctx context.Context) ([]Task, error)
	InProgress(ctx context.Context) ([]Task, error)
	Archived(ctx context.Context) ([]Task, error)
	WaitingForPayment(ctx context.Context) ([]Task, error)
}
