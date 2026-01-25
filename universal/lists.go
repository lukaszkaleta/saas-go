package universal

import "context"

type Lister[T Idable] interface {
	List(ctx context.Context) ([]T, error)
}

type FullText[T Idable] interface {
	ByQuery(ctx context.Context, query *string) ([]*T, error)
}
