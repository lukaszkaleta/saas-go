package universal

import "context"

// API

type Address interface {
	Model() *AddressModel
	Update(ctx context.Context, newModel *AddressModel) error
}

// Builder

// Model

type AddressModel struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	District   string `json:"district"`
}

func (model *AddressModel) Change(newModel *AddressModel) {
	model.Line1 = newModel.Line1
	model.Line2 = newModel.Line2
	model.City = newModel.City
	model.PostalCode = newModel.PostalCode
	model.District = newModel.District
}

func EmptyAddressModel() *AddressModel {
	return &AddressModel{
		Line1:      "",
		Line2:      "",
		City:       "",
		PostalCode: "",
		District:   "",
	}
}

// Solid

type SolidAddress struct {
	model   *AddressModel
	address Address
}

func NewSolidAddress(model *AddressModel, address Address) Address {
	return &SolidAddress{model, address}
}

func (addr SolidAddress) Update(ctx context.Context, newModel *AddressModel) error {
	addr.model.Change(newModel)
	if addr.address == nil {
		return nil
	}
	return addr.address.Update(ctx, newModel)
}

func (addr SolidAddress) Model() *AddressModel {
	return addr.model
}
