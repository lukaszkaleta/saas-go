package finance

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type SellerReporting interface {
	SumInPeriod(ctx context.Context, period universal.DateRange) (int64, error)
}
