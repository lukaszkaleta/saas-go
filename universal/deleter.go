package universal

import "context"

type Deleter interface {
	Delete(ctx context.Context) error
}
