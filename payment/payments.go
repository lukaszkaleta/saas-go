package payment

import "context"

type Payments interface {
	Create(ctx context.Context, offer any) (Intent, error)
	Search() Search
}
