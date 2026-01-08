package universal

import "context"

type Idable interface {
	ID() int64
}

type Idables[T Idable] interface {
	ById(ctx context.Context, id int64) (T, error)
}
