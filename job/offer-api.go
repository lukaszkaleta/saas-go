package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Offer interface {
	Model(ctx context.Context) (*OfferModel, error)
	Accept(ctx context.Context) error
	Reject(ctx context.Context) error
	Accepted() (bool, error)
	Rejected() (bool, error)
}

const Created = "created"
const Rejected = "rejected"
const Accepted = "accepted"

type OfferModel struct {
	Id          int64                       `json:"id"`
	JobId       int64                       `json:"jobId"`
	Price       *universal.PriceModel       `json:"price"`
	Description *universal.DescriptionModel `json:"description"`
	Rating      int                         `json:"rating"`
	Actions     universal.ActionsModel      `json:"actions"`
}

func EmptyOfferModel() *OfferModel {
	return &OfferModel{
		Id:          0,
		JobId:       0,
		Price:       universal.EmptyPriceModel(),
		Description: universal.EmptyDescriptionModel(),
		Rating:      0,
		Actions:     universal.ActionsModel{},
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

func (s *SolidOffer) Accept(ctx context.Context) error {
	if s.Offer != nil {
		err := s.Offer.Accept(ctx)
		if err != nil {
			return err
		}
	}
	now := time.Now()
	s.model.Actions.List[Accepted] = &universal.ActionModel{
		ById:   universal.CurrentUserId(ctx),
		MadeAt: &now,
		Name:   Accepted,
	}
	return nil
}

func (s *SolidOffer) Reject(ctx context.Context) error {
	if s.Offer != nil {
		err := s.Offer.Reject(ctx)
		if err != nil {
			return err
		}
	}
	now := time.Now()
	s.model.Actions.List[Rejected] = &universal.ActionModel{
		ById:   universal.CurrentUserId(ctx),
		MadeAt: &now,
		Name:   Rejected,
	}
	return nil
}

func (s *SolidOffer) Accepted() (bool, error) {
	actionModel := s.model.Actions.List[Accepted]
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (s *SolidOffer) Rejected() (bool, error) {
	actionModel := s.model.Actions.List[Rejected]
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (s *SolidOffer) Model(ctx context.Context) (*OfferModel, error) {
	return s.model, nil
}
