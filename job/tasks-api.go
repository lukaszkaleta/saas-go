package job

import "context"

type Tasks interface {
	Create(ctx context.Context, model *TaskModel) (Task, error)
	InProgress(ctx context.Context) ([]Task, error)
	Archived(ctx context.Context) ([]Task, error)
	WaitingForPayment(ctx context.Context) ([]Task, error)
}
