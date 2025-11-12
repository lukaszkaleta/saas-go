package job

import "github.com/lukaszkaleta/saas-go/universal"

type Offer interface {
	Accept(person universal.Person) error
	Reject(person universal.Person) error
}

type OfferModel struct {
	Id          int64                       `json:"id"`
	Price       *universal.PriceModel       `json:"price"`
	Description *universal.DescriptionModel `json:"description"`
	Rating      int                         `json:"rating"`
}

func EmptyOfferModel() *OfferModel {
	return &OfferModel{
		Id:          0,
		Price:       universal.EmptyPriceModel(),
		Description: universal.EmptyDescriptionModel(),
		Rating:      0,
	}
}

//
// Solid
//

func NewSolidOffer(model *OfferModel, offer Offer) Offer {
	return &SolidOffer{
		Id:    model.Id,
		model: model,
		Offer: offer,
	}
}

type SolidOffer struct {
	Id    int64
	model *OfferModel
	Offer Offer
}

func (s SolidOffer) Accept(person universal.Person) error {
	//TODO implement me
	panic("implement me")
}

func (s SolidOffer) Reject(person universal.Person) error {
	//TODO implement me
	panic("implement me")
}
