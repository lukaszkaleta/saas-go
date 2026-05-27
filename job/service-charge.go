package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type ServiceCharge interface {
	universal.Activable
	ID() int64
	Model(ctx context.Context) (*ServiceChargeModel, error)
}

type PriceMode string

const (
	PERCENT PriceMode = "PERCENT"
	FIXED   PriceMode = "FIXED"
)

type PriceFormula struct {
	Mode  PriceMode `json:"mode"`
	Value int       `json:"value"`
}

type ServiceChargeModel struct {
	Id     int64        `json:"id"`
	Worker PriceFormula `json:"worker"`
	Owner  PriceFormula `json:"owner"`
	Active bool         `json:"active"`
}

type SolidServiceCharge struct {
	model         *ServiceChargeModel
	serviceCharge ServiceCharge
}

func NewSolidServiceCharge(model *ServiceChargeModel, serviceCharge ServiceCharge) ServiceCharge {
	return &SolidServiceCharge{model, serviceCharge}
}

func (s *SolidServiceCharge) ID() int64 {
	return s.model.Id
}

func (s *SolidServiceCharge) Model(ctx context.Context) (*ServiceChargeModel, error) {
	return s.model, nil
}

func (s *SolidServiceCharge) Activate(ctx context.Context) error {
	return s.serviceCharge.Activate(ctx)
}

func (s *SolidServiceCharge) Deactivate(ctx context.Context) error {
	return s.serviceCharge.Deactivate(ctx)
}

func (s *SolidServiceCharge) IsActive(ctx context.Context) (bool, error) {
	return s.model.Active, nil
}
