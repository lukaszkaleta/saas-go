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

type ServiceChargeSummary struct {
	Charge       *ServiceChargeModel   `json:"charge"`
	Price        *universal.PriceModel `json:"price"`
	WorkerCharge *universal.PriceModel `json:"workerCharge"`
	WorkerCost   *universal.PriceModel `json:"workerCost"`
	OwnerCharge  *universal.PriceModel `json:"ownerCharge"`
	OwnerCost    *universal.PriceModel `json:"ownerCost"`
}

func NewServiceChargeSummary(charge *ServiceChargeModel, price *universal.PriceModel) *ServiceChargeSummary {
	c := &ServiceChargeSummary{
		Charge: charge,
		Price:  price,
	}
	c.WorkerCharge = c.calculate(charge.Worker, price)
	c.OwnerCharge = c.calculate(charge.Owner, price)
	c.WorkerCost = &universal.PriceModel{
		Value:    price.Value - c.WorkerCharge.Value,
		Currency: price.Currency,
	}
	c.OwnerCost = &universal.PriceModel{
		Value:    price.Value + c.OwnerCharge.Value,
		Currency: price.Currency,
	}
	return c
}

func (c *ServiceChargeSummary) WorkerPrice() *universal.PriceModel {
	return c.WorkerCost
}

func (c *ServiceChargeSummary) OwnerPrice() *universal.PriceModel {
	return c.OwnerCost
}

func (c *ServiceChargeSummary) calculate(formula PriceFormula, price *universal.PriceModel) *universal.PriceModel {
	val := 0
	switch formula.Mode {
	case PERCENT:
		val = price.Value * formula.Value / 100
	case FIXED:
		val = formula.Value
	}

	return &universal.PriceModel{
		Value:    val,
		Currency: price.Currency,
	}
}
