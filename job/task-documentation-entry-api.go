package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type TaskDocumentationEntry interface {
	universal.ModelAware[TaskDocumentationEntryModel]
	Images(ctx context.Context) ([]string, error)
	AddFirstImage(ctx context.Context, firstImageUrl string) error
}

type TaskDocumentationEntryModel struct {
	Id          int64                       `json:"id"`
	TaskId      int64                       `json:"taskId"`
	Summary     *universal.DescriptionModel `json:"summary"`
	Images      []string                    `json:"images"`
	CreatedById int64                       `json:"createdById"`
	CreatedAt   time.Time                   `json:"createdAt"`
}

func (t TaskDocumentationEntryModel) ID() int64 {
	return t.Id
}

type SolidTaskDocumentationEntry struct {
	model                  *TaskDocumentationEntryModel
	taskDocumentationEntry TaskDocumentationEntry
}

func NewSolidTaskDocumentationEntry(model *TaskDocumentationEntryModel, taskDocumentationEntry TaskDocumentationEntry) TaskDocumentationEntry {
	return &SolidTaskDocumentationEntry{model, taskDocumentationEntry}
}

func (solidEntry *SolidTaskDocumentationEntry) ID() int64 {
	return solidEntry.model.Id
}

func (solidEntry *SolidTaskDocumentationEntry) Model(ctx context.Context) (*TaskDocumentationEntryModel, error) {
	return solidEntry.model, nil
}

func (solidEntry *SolidTaskDocumentationEntry) Images(ctx context.Context) ([]string, error) {
	if solidEntry.taskDocumentationEntry != nil {
		return solidEntry.taskDocumentationEntry.Images(ctx)
	}
	return solidEntry.model.Images, nil
}

func (solidEntry *SolidTaskDocumentationEntry) AddFirstImage(ctx context.Context, firstImageUrl string) error {
	if solidEntry.taskDocumentationEntry != nil {
		return solidEntry.taskDocumentationEntry.AddFirstImage(ctx, firstImageUrl)
	}
	return nil
}
