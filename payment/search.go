package payment

import "context"

type Search interface {
	Intent(ctx context.Context, reference string) (Intent, error)
}
