package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type JobWorkers interface {
	Suggest(ctx context.Context, job Job) ([]*universal.PersonModel, error)
}
