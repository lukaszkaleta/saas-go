package payment

import "context"

type Payments interface {
	Creator() PaymentsCreator
	Search() Search
}

type PaymentsCreator interface {
	Create(ctx context.Context, offer any) (Intent, error)
}
