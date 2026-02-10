package universal

import "context"

type Idable interface {
	ID() int64
}

type JustId struct {
	id int64
}

func (ji JustId) ID() int64 {
	return ji.id
}

type Idables[T Idable] interface {
	ById(ctx context.Context, id int64) (T, error)
}
