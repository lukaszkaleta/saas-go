package universal

import "context"

type Creator[In any, Out any] interface {
	Create(ctx context.Context, in In) (Out, error)
}
