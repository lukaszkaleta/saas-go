package job

import (
	"context"
	"errors"

	"github.com/lukaszkaleta/saas-go/universal"
)

type TaskDocumentation interface {
	Create(ctx context.Context, summary *universal.DescriptionModel) (TaskDocumentationEntry, error)
	EntriesModels(ctx context.Context) ([]TaskDocumentationEntryModel, error)
}

var (
	ErrTaskDocumentationTaskNotFound = errors.New("task not found")
	ErrTaskDocumentationJobNotFound  = errors.New("job not found")
	ErrTaskDocumentationAccessDenied = errors.New("access denied")
	ErrTaskDocumentationMissingUser  = errors.New("missing current user")
)
