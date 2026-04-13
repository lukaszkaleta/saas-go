package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type TaskDocumentation interface {
	Create(ctx context.Context, summary *universal.DescriptionModel) (TaskDocumentationEntry, error)
	EntriesModels(ctx context.Context) ([]TaskDocumentationEntryModel, error)
}
