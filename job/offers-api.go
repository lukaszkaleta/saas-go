package job

import (
	"context"
	"fmt"

	"github.com/lukaszkaleta/saas-go/universal"
)

type OfferMaker interface {
	Make(ctx context.Context, workerId int64, model *OfferRevisionModel) (OfferRevision, error)
}

type OfferWaiter interface {
	Waiting(ctx context.Context) ([]Offer, error)
}

type Offers interface {
	universal.Deleter

	OfferMaker
	OfferWaiter

	ById(ctx context.Context, id int64) (Offer, error)
	FromUser(ctx context.Context, user universal.Idable) (Offer, error)
	Accepted(ctx context.Context) (Offer, error)
}

// No offers implementation

type NoOffers struct {
}

func (n NoOffers) ById(ctx context.Context, id int64) (Offer, error) {
	return nil, nil
}

func (n NoOffers) FromUser(ctx context.Context, user universal.Idable) (Offer, error) {
	return nil, nil
}
func (n NoOffers) Waiting(ctx context.Context) ([]Offer, error) {
	return nil, nil
}

func (n NoOffers) Make(ctx context.Context, workerId int64, model *OfferRevisionModel) (OfferRevision, error) {
	return nil, nil
}
func (n NoOffers) Accepted(ctx context.Context) (Offer, error) {
	return nil, nil
}

func (n NoOffers) Delete(ctx context.Context) error {
	return nil
}

type MessagesOfferMaker struct {
	inner OfferMaker
	job   Job
}

func NewMessagesOfferMaker(inner OfferMaker, job Job) *MessagesOfferMaker {
	return &MessagesOfferMaker{
		inner: inner,
		job:   job,
	}
}

func (m MessagesOfferMaker) Make(ctx context.Context, workerId int64, model *OfferRevisionModel) (OfferRevision, error) {
	revision, err := m.inner.Make(ctx, workerId, model)
	if err != nil {
		return nil, err
	}
	// We use model from argument which is OfferRevisionModel
	offerMessage := model.Description.Value
	// Make revision message a message which will be put into chat:
	message := fmt.Sprintf("%s: %s",
		model.Price.UserFriendly(),
		offerMessage,
	)

	revisionModel, err := revision.Model(ctx)
	if err != nil {
		return nil, err
	}
	userId := revisionModel.Actions.CreatedById()
	jobChat, err := m.job.Chats().Ensure(ctx, *userId)
	if err != nil {
		return nil, err
	}
	_, err = jobChat.Messages().AddGenerated(ctx, message, "offer")
	if err != nil {
		return nil, err
	}
	return revision, nil
}

func ModelsAndJobIds(ctx context.Context, list []Offer) ([]*OfferModel, []int64) {
	models := make([]*OfferModel, len(list))
	jobIds := make([]int64, len(list))
	for i, msg := range list {
		model, err := msg.Model(ctx)
		if err != nil {
			panic(err)
		}
		models[i] = model
		jobIds[i] = model.JobId
	}
	return models, jobIds
}

func Dtos(ctx context.Context, list []Offer) ([]*OfferDto, error) {
	dtos := make([]*OfferDto, len(list))
	for i, o := range list {
		model, err := o.Model(ctx)
		if err != nil {
			return nil, err
		}
		dto := &OfferDto{OfferModel: model}
		if model.LastRevisionId != nil {
			revision, err := o.Revisions().ById(ctx, *model.LastRevisionId)
			if err == nil && revision != nil {
				rm, err := revision.Model(ctx)
				if err == nil {
					dto.Revision = rm
				}
			}
		}
		dtos[i] = dto
	}
	return dtos, nil
}
