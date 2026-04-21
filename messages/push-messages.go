package messages

import (
	"context"
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

func (a *PushMessages) Add(ctx context.Context, model *MessageModel) (Message, error) {
	msg, err := a.messages.Add(ctx, model)

	// Do not send push on generated messages
	if model.ValueGenerated {
		return msg, err
	}

	if err != nil {
		return nil, err
	}

	recipient, err := a.users.ById(ctx, model.RecipientId)

	if err != nil {
		// We don't want to fail the message creation if push notification fails to start
		return msg, nil
	}

	account := recipient.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return msg, nil
	}

	body := model.Value
	if len(body) > 20 {
		body = body[:20] + "..."
	}

	jobId := strconv.FormatInt(model.OwnerId, 10)
	messageId := strconv.FormatInt(msg.ID(), 10)
	pushMsg := universal.PushMessage{
		Title: "New Message",
		Body:  body,
		Link:  "https://naborly.no/chat" + string(jobId) + "/" + messageId,
	}

	a.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return msg, nil
}

func (a *PushMessages) List(ctx context.Context) ([]Message, error) {
	return a.messages.List(ctx)
}

func (a *PushMessages) ById(ctx context.Context, id int64) (Message, error) {
	return a.messages.ById(ctx, id)
}

func (a *PushMessages) AddSimple(ctx context.Context, recipientId int64, value string) (Message, error) {
	return a.messages.AddSimple(ctx, recipientId, value)
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
