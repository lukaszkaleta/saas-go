package universal

import "fmt"

// API

type Price interface {
	Model() *PriceModel
	Update(newModel *PriceModel) error
}

// Builder

func PriceFromModel(model *PriceModel) Price {
	return SolidPrice{
		model: model,
	}
}

// Model

type PriceModel struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

func EmptyPriceModel() *PriceModel {
	return &PriceModel{
		Value:    0,
		Currency: "",
	}
}

func (model *PriceModel) Change(newModel *PriceModel) {
	model.Value = newModel.Value
	model.Currency = newModel.Currency
}

func (model *PriceModel) UserFriendly() string {
	return fmt.Sprintf("%s %s", model.DecimalValue(), model.Currency)
}

func (model *PriceModel) DecimalValue() string {
	return fmt.Sprintf("%.2f", float64(model.Value)/100)
}

// Solid

type SolidPrice struct {
	model *PriceModel
	price Price
}

func NewSolidPrice(model *PriceModel, Price Price) Price {
	return &SolidPrice{model, Price}
}

func (addr SolidPrice) Update(newModel *PriceModel) error {
	addr.model.Change(newModel)
	if addr.price == nil {
		return nil
	}
	return addr.price.Update(newModel)
}

func (addr SolidPrice) Model() *PriceModel {
	return addr.model
}
