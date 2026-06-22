package job

import (
	"context"
	"log/slog"
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
	universal.Idable `json:"-"`
	Id               int64                       `json:"id"`
	OfferId          int64                       `json:"offerId"`
	Price            *universal.PriceModel       `json:"price"`
	Description      *universal.DescriptionModel `json:"description"`
	Actions          *universal.ActionsModel     `json:"actions"`
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

//
// When rejecting offer revision we need to generate message
//

type MessagesOfferRejecter struct {
	inner OfferRevision
	offer Offer
	job   Job
}

func (m *MessagesOfferRejecter) Reject(ctx context.Context) error {
	err := m.inner.Reject(ctx)
	if err != nil {
		return err
	}

	// Check who created offer
	userId, err := universal.CreatedById[OfferRevisionModel](ctx, m.inner)
	if err != nil {
		slog.Error("Failed to get creator ID for message", "error", err)
		return nil // Main operation succeeded
	}
	// AddSimple message that offer is rejected
	jobChat, err := m.job.Chats().Ensure(ctx, userId)
	if err != nil {
		slog.Error("Failed to get chat for job", "error", err)
		return nil
	}
	_, err = jobChat.Messages().AddGenerated(ctx, "Offer rejected", "offer rejected")
	if err != nil {
		slog.Error("Failed to add message", "error", err)
		return nil // Main operation succeeded
	}
	return nil
}

func NewMessagesOfferRejecter(job Job, offer Offer, inner OfferRevision) universal.Rejecter {
	return &MessagesOfferRejecter{
		inner: inner,
		offer: offer,
		job:   job,
	}
}

//
// When accepting offer we need to send a message
//

type MessagesOfferRevisionAcceptor struct {
	inner universal.Acceptor
	job   Job
}

func (m *MessagesOfferRevisionAcceptor) Accept(ctx context.Context) error {
	err := m.inner.Accept(ctx)
	if err != nil {
		return err
	}

	// Check who created offer
	userId, err := universal.CreatedById[OfferRevisionModel](ctx, m.inner)
	if err != nil {
		slog.Error("Failed to get creator ID for message", "error", err)
		return nil // Main operation succeeded
	}
	jobChat, err := m.job.Chats().Ensure(ctx, userId)
	if err != nil {
		slog.Error("Failed to get chat for job", "error", err)
		return err
	}
	_, err = jobChat.Messages().AddGenerated(ctx, "Offer accepted", "offer accepted")
	if err != nil {
		slog.Error("Failed to add message", "error", err)
		return nil // Main operation succeeded
	}
	return nil
}

func NewMessagesOfferRevisionAcceptor(job Job, inner universal.Acceptor) universal.Acceptor {
	return &MessagesOfferRevisionAcceptor{
		inner: inner,
		job:   job,
	}
}

//
// When accepting offer job will be moved to Occupied state
//

type ApproveOfferRevisionAcceptor struct {
	inner universal.Acceptor
	job   Job
}

func (m *ApproveOfferRevisionAcceptor) Accept(ctx context.Context) error {
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
	return &ApproveOfferRevisionAcceptor{job: job, inner: inner}
}
