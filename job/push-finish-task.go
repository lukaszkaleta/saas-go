package job

import (
	"context"
	"log/slog"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushFinishTask struct {
	task   Task
	users  user.Users
	sender *universal.PushSender
}

func NewPushTask(task Task, users user.Users) *PushFinishTask {
	return &PushFinishTask{
		task:   task,
		users:  users,
		sender: &universal.PushSender{},
	}
}

func (t *PushFinishTask) Finish(ctx context.Context) error {
	err := t.task.Finish(ctx)
	if err != nil {
		return err
	}

	pushErr := t.sendPush(ctx)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return nil
}

func (t *PushFinishTask) sendPush(ctx context.Context) error {
	model, err := t.task.Model(ctx)
	if err != nil {
		return err
	}

	job, err := t.task.Job(ctx)
	if err != nil {
		return err
	}
	jobModel, err := job.Model(ctx)
	if err != nil {
		return err
	}
	recipientId := jobModel.Actions.CreatedById()
	recipient, err := t.users.ById(ctx, *recipientId)
	if err != nil {
		return err
	}

	account := recipient.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return nil
	}

	taskUser, err := t.users.ById(ctx, model.UserId)
	if err != nil {
		return err
	}
	personModel := taskUser.Person().Model(ctx)

	pushMsg := universal.PushMessage{
		Title: personModel.FirstName,
		Body:  "Have finished your job, please review it",
		Link:  "https://naborly.no/owner/dashboard",
	}

	t.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
