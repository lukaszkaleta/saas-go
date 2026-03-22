package payment

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Intent interface {
	universal.Idable
	Model(ctx context.Context) (*IntentModel, error)
}

type IntentModel struct {
	Id                    int64                   `json:"id"`
	Reference             string                  `json:"reference"`
	StripePaymentIntentId string                  `json:"stripePaymentIntentId,omitzero"`
	StripeClientSecret    string                  `json:"stripeClientSecret,omitzero"`
	JobId                 int64                   `json:"jobId"`
	PayerId               int64                   `json:"payerId"`
	PayeeId               int64                   `json:"payeeId"`
	Amount                int64                   `json:"amount"`
	Currency              string                  `json:"currency"`
	Status                string                  `json:"status"`
	Actions               *universal.ActionsModel `json:"actions"`
}

func (m IntentModel) ID() int64 {
	return m.Id
}

func (m IntentModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

// Solid

type SolidIntent struct {
	Id     int64
	model  *IntentModel
	intent Intent
}

func NewSolidIntent(model *IntentModel, intent Intent, id int64) Intent {
	return &SolidIntent{
		Id:     id,
		model:  model,
		intent: intent,
	}
}

func (m *SolidIntent) Model(ctx context.Context) (*IntentModel, error) {
	return m.model, nil
}

func (m *SolidIntent) ID() int64 {
	return m.Id
}

func EmptyIntentModel() *IntentModel {
	return &IntentModel{
		Actions: universal.EmptyActionsModel(),
	}
}
