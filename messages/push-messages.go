package messages

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushMessages struct {
	messages Messages
	users    user.Users
	sender   *universal.PushSender
}

func NewPushMessages(messages Messages, users user.Users, sender *universal.PushSender) *PushMessages {
	return &PushMessages{
		messages: messages,
		users:    users,
		sender:   sender,
	}
}

func (a *PushMessages) sendPush(ctx context.Context, msg Message) error {
	model, err := msg.Model(ctx)
	if err != nil {
		return err
	}

	recipient, err := a.users.ById(ctx, model.RecipientId)
	if err != nil {
		// We don't want to fail the message creation if push notification fails to start
		return err
	}

	account := recipient.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return nil
	}

	body := model.Value
	if len(body) > 20 {
		body = body[:20] + "..."
	}

	jobId := strconv.FormatInt(model.OwnerId, 10)
	messageId := strconv.FormatInt(msg.ID(), 10)
	pushMsg := universal.PushMessage{
		Title: user.CurrentUser(ctx).Person.FirstName,
		Body:  body,
		Link:  "https://naborly.no/chat/" + string(jobId) + "/" + messageId,
	}

	a.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}

func (a *PushMessages) Add(ctx context.Context, model *MessageModel) (Message, error) {
	return a.messages.Add(ctx, model)
}

func (a *PushMessages) List(ctx context.Context) ([]Message, error) {
	return a.messages.List(ctx)
}

func (a *PushMessages) ById(ctx context.Context, id int64) (Message, error) {
	return a.messages.ById(ctx, id)
}

func (a *PushMessages) AddSimple(ctx context.Context, recipientId int64, value string) (Message, error) {
	msg, err := a.messages.AddSimple(ctx, recipientId, value)
	pushError := a.sendPush(ctx, msg)
	if pushError != nil {
		slog.Error("Can not send push", "Error", pushError.Error())
	}
	return msg, err
}

func (a *PushMessages) AddGenerated(ctx context.Context, recipientId int64, value string) (Message, error) {
	return a.messages.AddGenerated(ctx, recipientId, value)
}

func (a *PushMessages) ForRecipient(ctx context.Context, recipient universal.Idable) ([]Message, error) {
	return a.messages.ForRecipient(ctx, recipient)
}

func (a *PushMessages) ForRecipientById(ctx context.Context, id int64) ([]Message, error) {
	return a.messages.ForRecipientById(ctx, id)
}

func (a *PushMessages) Acknowledge(ctx context.Context) error {
	return a.messages.Acknowledge(ctx)
}

func (a *PushMessages) LastQuestions(ctx context.Context) ([]Message, error) {
	return a.messages.LastQuestions(ctx)
}
