package universal

import "context"

type Adder[In any, Out any] interface {
	Add(ctx context.Context, in In) (Out, error)
}
