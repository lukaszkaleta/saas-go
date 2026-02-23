package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type Task interface {
	universal.Idable
	filestore.FileSystemAware
	Model(ctx context.Context) (*TaskModel, error)
	Description() universal.Description
}

type TaskModel struct {
	Id          int64                       `json:"id"`
	JobId       int64                       `json:"jobId"`
	UserId      int64                       `json:"userId"`
	OfferId     int64                       `json:"offerId"`
	Description *universal.DescriptionModel `json:"description"`
	Actions     *universal.ActionsModel     `json:"actions"`
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

func (s *SolidTask) Description() universal.Description {
	if s.Task != nil {
		return universal.NewSolidDescription(
			s.model.Description,
			s.Task.Description(),
		)
	}
	return universal.NewSolidDescription(s.model.Description, nil)
}

func (s *SolidTask) FileSystem() filestore.FileSystem {
	return s.Task.FileSystem()
}
