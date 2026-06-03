package universal

import "context"

type Counter interface {
	Increment(ctx context.Context) error
	Get(ctx context.Context) (int64, error)
}
