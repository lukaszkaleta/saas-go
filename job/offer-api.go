package job

import (
	"context"
	"log/slog"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Offer interface {
	universal.Idable
	universal.ModelAware[OfferModel]
	universal.Creator[int64, Offer]
	Revisions() OfferRevisions
	Accepted() (bool, error)
	Rejected() (bool, error)
	Accept(ctx context.Context) error
	Reject(ctx context.Context) error
}

const Created = "created"
const Rejected = "rejected"
const Accepted = "accepted"

type OfferDto struct {
	*OfferModel
	Revision *OfferRevisionModel `json:"revision"`
}

type OfferModel struct {
	universal.Idable   `json:"-"`
	Id                 int64  `json:"id"`
	JobId              int64  `json:"jobId"`
	WorkerId           int64  `json:"workerId"`
	AcceptedRevisionId *int64 `json:"acceptedRevisionId"`
	LastRevisionId     *int64 `json:"lastRevisionId"`
	Status             string `json:"status"`
	Rating             int    `json:"rating"`
}

func EmptyOfferModel() *OfferModel {
	om := &OfferModel{
		Id:    0,
		JobId: 0,
	}
	return om
}

//
// Solid
//

func (s *SolidOffer) Create(ctx context.Context, in int64) (Offer, error) {
	return s.Offer.Create(ctx, in)
}

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

func (s *SolidOffer) Model(ctx context.Context) (*OfferModel, error) {
	return s.model, nil
}

func (s *SolidOffer) Revisions() OfferRevisions {
	return s.Offer.Revisions()
}

func (s *SolidOffer) Accepted() (bool, error) {
	return s.Offer.Accepted()
}

func (s *SolidOffer) Rejected() (bool, error) {
	return s.Offer.Rejected()
}

func (s *SolidOffer) Accept(ctx context.Context) error {
	return s.Offer.Accept(ctx)
}

func (s *SolidOffer) Reject(ctx context.Context) error {
	return s.Offer.Reject(ctx)
}

//
// When accepting offer we will create a task for user who created offer.
//

type TaskOnOfferAccept struct {
	inner   universal.Acceptor
	offerId int64
	job     Job
}

func (m *TaskOnOfferAccept) Accept(ctx context.Context) error {
	err := m.inner.Accept(ctx)
	if err != nil {
		return err
	}
	err = m.job.MakeTask(ctx, m.offerId)
	if err != nil {
		slog.Error("Failed to make task", "error", err)
		return nil // Main operation succeeded
	}

	return nil
}

func NewTaskOnOfferAccept(job Job, offerId int64, inner universal.Acceptor) universal.Acceptor {
	return &TaskOnOfferAccept{
		inner:   inner,
		offerId: offerId,
		job:     job,
	}
}
