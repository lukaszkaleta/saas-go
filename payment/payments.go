package payment

import "context"

type Payments interface {
	Create(ctx context.Context, id int64) (Intent, error)
	Search() Search
}
