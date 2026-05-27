package job

import (
	"context"
)

type ServiceCharges interface {
	Active(ctx context.Context) (ServiceCharge, error)
	Add(ctx context.Context, model *ServiceChargeModel) (ServiceCharge, error)
	All(ctx context.Context) ([]ServiceCharge, error)
}
