package universal

import "context"

type Lister[T Idable] interface {
	List(ctx context.Context) ([]T, error)
}
