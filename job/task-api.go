package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type Task interface {
	universal.Idable
	filestore.FileSystemAware
	universal.ActionsAware
	Model(ctx context.Context) (*TaskModel, error)
	Summary() universal.Description
	Job(ctx context.Context) (Job, error)
	Finish(ctx context.Context) error
	Documentation() TaskDocumentation
}

type TaskModel struct {
	Id      int64                       `json:"id"`
	JobId   int64                       `json:"jobId"`
	UserId  int64                       `json:"userId"`
	OfferId int64                       `json:"offerId"`
	Summary *universal.DescriptionModel `json:"summary"`
	Actions *universal.ActionsModel     `json:"actions"`
}

func (m TaskModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

// Solid

func NewSolidTask(model *TaskModel, task Task) Task {
	return &SolidTask{
		Id:    model.Id,
		model: model,
		Task:  task,
	}
}

type SolidTask struct {
	Id    int64
	model *TaskModel
	Task  Task
}

func (s *SolidTask) ID() int64 {
	return s.Id
}

func (s *SolidTask) Model(ctx context.Context) (*TaskModel, error) {
	return s.model, nil
}

func (s *SolidTask) Summary() universal.Description {
	if s.Task != nil {
		return universal.NewSolidDescription(
			s.model.Summary,
			s.Task.Summary(),
		)
	}
	return universal.NewSolidDescription(s.model.Summary, nil)
}

func (s *SolidTask) FileSystem() filestore.FileSystem {
	return s.Task.FileSystem()
}

func (s *SolidTask) Job(ctx context.Context) (Job, error) {
	return s.Task.Job(ctx)
}

func (s *SolidTask) Finish(ctx context.Context) error {
	err := s.Task.Finish(ctx)
	if err != nil {
		return err
	}
	model, err := s.Model(ctx)
	if err != nil {
		return err
	}
	finishedAction := model.Actions.WithName("finished")
	if finishedAction != nil {
		now := time.Now()
		finishedAction.MadeAt = &now
		finishedAction.ById = universal.CurrentUserId(ctx)
	}
	return nil
}

func (s *SolidTask) Actions() universal.Actions {
	return s.Task.Actions()
}

func (s *SolidTask) Documentation() TaskDocumentation {
	return s.Task.Documentation()
}
