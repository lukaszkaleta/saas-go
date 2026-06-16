package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type OfferRevision interface {
	universal.Idable
	universal.Acceptable
	universal.Rejectable
	universal.ActionsAware
	universal.ModelAware[OfferRevisionModel]
}

type OfferRevisionModel struct {
	universal.Idable
	Id          int64                       `json:"id"`
	OfferId     int64                       `json:"offerId"`
	Price       *universal.PriceModel       `json:"price"`
	Description *universal.DescriptionModel `json:"description"`
	Actions     *universal.ActionsModel     `json:"actions"`
}

func (m OfferRevisionModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

func EmptyOfferRevisionModel() *OfferRevisionModel {
	orm := &OfferRevisionModel{
		Id:          0,
		OfferId:     0,
		Price:       universal.EmptyPriceModel(),
		Description: universal.EmptyDescriptionModel(),
	}
	orm.Actions = universal.EmptyActionsModel()
	return orm
}

//
// Solid
//

func NewSolidOfferRevision(model *OfferRevisionModel, revision OfferRevision) OfferRevision {
	return &SolidOfferRevision{
		Id:            model.Id,
		model:         model,
		OfferRevision: revision,
	}
}

type SolidOfferRevision struct {
	universal.Idable
	Id            int64
	model         *OfferRevisionModel
	OfferRevision OfferRevision
}

func (s *SolidOfferRevision) ID() int64 {
	return s.Id
}

func (s *SolidOfferRevision) Accept(ctx context.Context) error {
	if s.OfferRevision != nil {
		err := s.OfferRevision.Accept(ctx)
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

func (s *SolidOfferRevision) Accepted() (bool, error) {
	actionModel := s.model.Actions.List[Accepted]
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (s *SolidOfferRevision) Reject(ctx context.Context) error {
	if s.OfferRevision != nil {
		err := s.OfferRevision.Reject(ctx)
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

func (s *SolidOfferRevision) Rejected() (bool, error) {
	actionModel := s.model.Actions.List[Rejected]
	if actionModel == nil {
		return false, nil
	}
	return actionModel.Exists(), nil
}

func (s *SolidOfferRevision) Model(ctx context.Context) (*OfferRevisionModel, error) {
	return s.model, nil
}

func (s *SolidOfferRevision) Actions() universal.Actions {
	if s.OfferRevision != nil {
		return s.OfferRevision.Actions()
	}
	return nil
}
