package job

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

// Push To Job Owner

type PushToJobOwnerMakeOffer struct {
	inner  OfferMaker
	job    Job
	users  user.Users
	sender *universal.PushSender
}

func NewPushToJobOwnerMakeOffer(inner OfferMaker, job Job, users user.Users) *PushToJobOwnerMakeOffer {
	return &PushToJobOwnerMakeOffer{
		inner:  inner,
		job:    job,
		users:  users,
		sender: &universal.PushSender{},
	}
}

func (p *PushToJobOwnerMakeOffer) Make(ctx context.Context, workerId int64, model *OfferRevisionModel) (OfferRevision, error) {
	revision, err := p.inner.Make(ctx, workerId, model)
	if err != nil {
		return nil, err
	}

	pushErr := p.sendPush(ctx, revision)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return revision, nil
}

func (p *PushToJobOwnerMakeOffer) sendPush(ctx context.Context, revision OfferRevision) error {
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
	personModel := sender.Person(ctx).Model(ctx)

	jobId := strconv.FormatInt(jobModel.ID(), 10)
	revisionModel, err := revision.Model(ctx)
	offerId := strconv.FormatInt(revisionModel.OfferId, 10)
	revisionId := strconv.FormatInt(revision.ID(), 10)

	pushMsg := universal.PushMessage{
		Title: personModel.FirstName,
		Body:  "Send you an offer",
		Link:  "https://naborly.no/offer/" + jobId + "/" + offerId + "/" + revisionId,
	}

	p.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}

// Push to Worker

type PushToWorkerMakeOffer struct {
	inner  OfferMaker
	job    Job
	users  user.Users
	sender *universal.PushSender
}

func NewPushToWorkerMakeOffer(inner OfferMaker, job Job, users user.Users) *PushToWorkerMakeOffer {
	return &PushToWorkerMakeOffer{
		inner:  inner,
		job:    job,
		users:  users,
		sender: &universal.PushSender{},
	}
}

func (p *PushToWorkerMakeOffer) Make(ctx context.Context, workerId int64, model *OfferRevisionModel) (OfferRevision, error) {
	revision, err := p.inner.Make(ctx, workerId, model)
	if err != nil {
		return nil, err
	}

	pushErr := p.sendPush(ctx, revision)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return revision, nil
}

func (p *PushToWorkerMakeOffer) sendPush(ctx context.Context, revision OfferRevision) error {

	revisionModel, err := revision.Model(ctx)
	if err != nil {
		return err
	}

	workerId := revisionModel.Actions.CreatedById()
	if workerId == nil {
		return nil
	}

	recipient, err := p.users.ById(ctx, *workerId)
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
	personModel := sender.Person(ctx).Model(ctx)

	jobId := strconv.FormatInt(p.job.ID(), 10)
	offerId := strconv.FormatInt(revisionModel.OfferId, 10)
	revisionId := strconv.FormatInt(revision.ID(), 10)

	pushMsg := universal.PushMessage{
		Title: personModel.FirstName,
		Body:  "Send you an offer",
		Link:  "https://naborly.no/offer/" + jobId + "/" + offerId + "/" + revisionId,
	}

	p.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
