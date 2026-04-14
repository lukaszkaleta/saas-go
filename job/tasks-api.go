package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Tasks interface {
	Create(ctx context.Context, model *TaskModel) (Task, error)
	ById(ctx context.Context, id int64) (Task, error)
	ByJobId(ctx context.Context, jobId int64) (Task, error)
	Current(ctx context.Context) (*TasksResult, error)
	InProgress(ctx context.Context) ([]Task, error)
	Completed(ctx context.Context) (*TasksResult, error)
	WaitingForPayment(ctx context.Context) ([]Task, error)
	Earnings(ctx context.Context) (map[string]*universal.PriceModel, error)
}

type TasksResult struct {
	Tasks  []Task
	Jobs   map[int64]Job
	People map[int64]*universal.PersonModel
}

func (tsResult *TasksResult) Job(ctx context.Context, task Task) (Job, error) {
	model, err := task.Model(ctx)
	if err != nil {
		return nil, err
	}
	job := tsResult.Jobs[model.JobId]
	return job, nil
}

func (tsResult *TasksResult) PersonModel(ctx context.Context, job Job) (*universal.PersonModel, error) {
	model, err := job.Model(ctx)
	if err != nil {
		return nil, err
	}
	personModel := tsResult.People[*model.Actions.CreatedById()]
	return personModel, nil
}

func NewTasksResult(tasks []Task, jobs []Job, people []*universal.PersonModel) *TasksResult {
	return &TasksResult{
		Tasks:  tasks,
		Jobs:   universal.IdableToMap(jobs),
		People: universal.IdableToMap(people),
	}
}
