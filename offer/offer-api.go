package offer

import (
	"time"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Offer interface {
	Model() *OfferModel
	Address() universal.Address
	Position() universal.Position
	Price() universal.Price
	Description() universal.Description
	FileSystem() filestore.FileSystem
	State() universal.State
}

type OfferStatus struct {
	Draft     time.Time `json:"draft"`
	Published time.Time `json:"published"`
	Closed    time.Time `json:"closed"`
}

func (o *OfferStatus) Current() string {
	if !o.Closed.IsZero() {
		return OfferClosed
	}
	if !o.Published.IsZero() {
		return OfferClosed
	}
	return OfferDraft
}

const (
	OfferDraft     string = "draft"
	OfferPublished string = "published"
	OfferClosed    string = "closed"
)

func OfferStatuses() []string {
	return []string{OfferDraft, OfferPublished, OfferClosed}
}

// Model

type OfferModel struct {
	Id          int64                       `json:"id"`
	Position    *universal.PositionModel    `json:"position"`
	Price       *universal.PriceModel       `json:"price"`
	Address     *universal.AddressModel     `json:"address"`
	Description *universal.DescriptionModel `json:"description"`
	State       OfferStatus                 `json:"state"`
}

func (m *OfferModel) Hint() *OfferHint {
	return &OfferHint{
		Id:    m.Id,
		Price: m.Price.UserFriendly(),
	}
}

func EmptyOfferModel() *OfferModel {
	return &OfferModel{
		Id:          0,
		Position:    &universal.PositionModel{},
		Address:     &universal.AddressModel{},
		Description: universal.EmptyDescriptionModel(),
		Price:       &universal.PriceModel{},
	}
}

type OfferHint struct {
	Id       int64                    `json:"id"`
	Position *universal.PositionModel `json:"position"`
	Price    string                   `json:"price"`
}

// Solid

type SolidOffer struct {
	Id    int64
	model *OfferModel
	Offer Offer
}

func NewSolidOffer(model *OfferModel, Offer Offer, id int64) Offer {
	return &SolidOffer{
		id,
		model,
		Offer,
	}
}
func (solidOffer *SolidOffer) Model() *OfferModel {
	return solidOffer.model
}

func (solidOffer *SolidOffer) Position() universal.Position {
	if solidOffer.Offer != nil {
		return universal.NewSolidPosition(
			solidOffer.Model().Position,
			solidOffer.Offer.Position(),
		)
	}
	return universal.NewSolidPosition(solidOffer.Model().Position, nil)
}

func (solidOffer *SolidOffer) Price() universal.Price {
	if solidOffer.Offer != nil {
		return universal.NewSolidPrice(
			solidOffer.Model().Price,
			solidOffer.Offer.Price(),
		)
	}
	return universal.NewSolidPrice(solidOffer.Model().Price, nil)
}

func (solidOffer *SolidOffer) Address() universal.Address {
	if solidOffer.Offer != nil {
		return universal.NewSolidAddress(
			solidOffer.Model().Address,
			solidOffer.Offer.Address(),
		)
	}
	return universal.NewSolidAddress(solidOffer.Model().Address, nil)
}
func (solidOffer *SolidOffer) Description() universal.Description {
	if solidOffer.Offer != nil {
		return universal.NewSolidDescription(
			solidOffer.Model().Description,
			solidOffer.Offer.Description(),
		)
	}
	return universal.NewSolidDescription(solidOffer.Model().Description, nil)
}

func (solidOffer *SolidOffer) FileSystem() filestore.FileSystem {
	return solidOffer.Offer.FileSystem()
}

func (solidOffer *SolidOffer) State() universal.State {
	available := OfferStatuses()
	current := solidOffer.Model().State.Current()
	if solidOffer.Offer != nil {
		return universal.NewSolidState(
			current,
			available,
			solidOffer.Offer.State())
	}
	return universal.NewSolidState(current, available, nil)
}
