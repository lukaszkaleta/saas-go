package job

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushMakeOffer struct {
	inner  OfferMaker
	job    Job
	users  user.Users
	sender *universal.PushSender
}

func NewPushMakeOffer(inner OfferMaker, job Job, users user.Users) *PushMakeOffer {
	return &PushMakeOffer{
		inner:  inner,
		job:    job,
		users:  users,
		sender: &universal.PushSender{},
	}
}

func (p *PushMakeOffer) Make(ctx context.Context, model *OfferModel) (Offer, error) {
	offer, err := p.inner.Make(ctx, model)
	if err != nil {
		return nil, err
	}

	pushErr := p.sendPush(ctx, offer)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return offer, nil
}

func (p *PushMakeOffer) sendPush(ctx context.Context, offer Offer) error {
	jobModel, err := p.job.Model(ctx)
	if err != nil {
		return err
	}

	ownerId := jobModel.Actions.CreatedById()
	if ownerId == nil {
		return nil
	}

	recipient, err := p.users.ById(ctx, *ownerId)
	if err != nil {
		return err
	}

	account := recipient.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return nil
	}

	senderId := universal.CurrentUserId(ctx)
	sender, err := p.users.ById(ctx, *senderId)
	if err != nil {
		return err
	}
	personModel := sender.Person().Model(ctx)

	jobId := strconv.FormatInt(jobModel.ID(), 10)
	offerId := strconv.FormatInt(offer.ID(), 10)

	pushMsg := universal.PushMessage{
		Title: personModel.FirstName,
		Body:  "Send you an offer",
		Link:  "https://naborly.no/offer/" + jobId + "/" + offerId,
	}

	p.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
