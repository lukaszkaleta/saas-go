package payment

import "context"

type Payments interface {
	Create(ctx context.Context) (Intent, error)
	Search() Search
}
