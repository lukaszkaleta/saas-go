package job

import (
	"context"
	"log/slog"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushCloserJob struct {
	job    Job
	users  user.Users
	sender *universal.PushSender
}

func NewPushCloserJob(job Job, users user.Users) *PushCloserJob {
	return &PushCloserJob{
		job:    job,
		users:  users,
		sender: &universal.PushSender{},
	}
}

func (t *PushCloserJob) Close(ctx context.Context) error {
	err := t.job.Close(ctx)
	if err != nil {
		return err
	}

	pushErr := t.sendPush(ctx)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return nil
}

func (t *PushCloserJob) sendPush(ctx context.Context) error {
	jobModel, err := t.job.Model(ctx)
	if err != nil {
		return err
	}

	offers := t.job.Offers()
	acceptedOffer, err := offers.Accepted(ctx)
	if err != nil || acceptedOffer == nil {
		return err
	}

	offerModel, err := acceptedOffer.Model(ctx)
	if err != nil {
		return err
	}

	userId := offerModel.Actions.CreatedById()
	if userId == nil {
		return nil
	}

	recipient, err := t.users.ById(ctx, *userId)
	if err != nil {
		return err
	}

	account := recipient.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return nil
	}

	owner, err := t.users.ById(ctx, *jobModel.Actions.CreatedById())
	if err != nil {
		return err
	}
	personModel := owner.Person().Model(ctx)

	pushMsg := universal.PushMessage{
		Title: personModel.FirstName,
		Body:  "Approved your work",
		Link:  "https://naborly.no/worker/dashboard",
	}

	t.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
