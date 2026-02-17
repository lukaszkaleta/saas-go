package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Offer interface {
	universal.Idable
	universal.Acceptable
	universal.Rejectable
	universal.ActionsAware
	universal.ModelAware[OfferModel]
}

const Created = "created"
const Rejected = "rejected"
const Accepted = "accepted"

type OfferModel struct {
	universal.Idable
	Id          int64                       `json:"id"`
	JobId       int64                       `json:"jobId"`
	Price       *universal.PriceModel       `json:"price"`
	Description *universal.DescriptionModel `json:"description"`
	Rating      int                         `json:"rating"`
	Actions     *universal.ActionsModel     `json:"actions"`
}

func (m OfferModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

func EmptyOfferModel() *OfferModel {
	om := &OfferModel{
		Id:          0,
		JobId:       0,
		Price:       universal.EmptyPriceModel(),
		Description: universal.EmptyDescriptionModel(),
		Rating:      0,
	}
	om.Actions = universal.EmptyActionsModel()
	return om
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
	universal.Idable
	Id    int64
	model *OfferModel
	Offer Offer
}

func (s *SolidOffer) ID() int64 {
	return s.Id
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

func (s *SolidOffer) Actions() universal.Actions {
	return s.Offer.Actions()
}

//
// When accepting offer we need to send a message
//

type MessagesOfferAcceptor struct {
	inner Offer
	job   Job
}

func (m *MessagesOfferAcceptor) Accept(ctx context.Context) error {
	userId, err := universal.CreatedById[JobModel](ctx, m.job)
	if err != nil {
		return err
	}
	_, err = m.job.Messages().Add(ctx, userId, "Offer accepted")
	if err != nil {
		return err
	}
	return m.inner.Accept(ctx)
}

//
// When accepting offer job will be moved to Occupied state
//

func NewMessagesOfferAcceptor(job Job, inner Offer) universal.Acceptor {
	return &MessagesOfferAcceptor{
		inner: inner,
		job:   job,
	}
}

type ApproveOfferAcceptor struct {
	inner universal.Acceptor
	job   Job
}

func (m *ApproveOfferAcceptor) Accept(ctx context.Context) error {
	err := m.inner.Accept(ctx)
	if err != nil {
		return err
	}
	err = m.job.State().Change(ctx, JobOccupied)
	if err != nil {
		return err
	}
	return nil
}

func NewApproveOfferAcceptor(job Job, inner universal.Acceptor) universal.Acceptor {
	return &ApproveOfferAcceptor{
		inner: inner,
		job:   job,
	}
}

//
// When rejecting offer
//

type MessagesOfferRejecter struct {
	inner Offer
	job   Job
}

func (m *MessagesOfferRejecter) Reject(ctx context.Context) error {
	// Check who created offer
	offerModel, err := m.inner.Model(ctx)
	if err != nil {
		return err
	}
	userId := offerModel.Actions.CreatedById()
	// Add message that offer is rejected
	_, err = m.job.Messages().Add(ctx, *userId, "Offer rejected")
	if err != nil {
		return err
	}
	// Reject offer.
	return m.inner.Reject(ctx)
}

func NewMessagesOfferRejecter(job Job, inner Offer) universal.Rejecter {
	return &MessagesOfferRejecter{
		inner: inner,
		job:   job,
	}
}
