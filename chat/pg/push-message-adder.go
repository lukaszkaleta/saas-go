package pg

import (
	"context"
	"strconv"

	"github.com/lukaszkaleta/saas-go/chat"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushOnMessageCreator struct {
	universal.Creator[string, chat.Message]

	inner     universal.Creator[string, chat.Message]
	recipient user.User
	chat      chat.Chat
	sender    *universal.PushSender
}

func NewPushOnMessageCreator(inner universal.Creator[string, chat.Message], recipient user.User, chat chat.Chat, sender *universal.PushSender) *PushOnMessageCreator {
	return &PushOnMessageCreator{
		inner:     inner,
		recipient: recipient,
		chat:      chat,
		sender:    sender,
	}
}

func (a *PushOnMessageCreator) Create(ctx context.Context, in string) (chat.Message, error) {
	msg, err := a.inner.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	_ = a.sendPush(ctx, msg)
	return msg, nil
}

func (a *PushOnMessageCreator) sendPush(ctx context.Context, msg chat.Message) error {
	model, err := msg.Model(ctx)
	if err != nil {
		return err
	}

	if a.recipient == nil {
		return nil
	}
	account := a.recipient.Account()
	if account == nil {
		return nil
	}
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel == nil || accountModel.FirebaseToken == "" {
		return nil
	}

	chatModel, err := a.chat.Model(ctx)
	if err != nil {
		return err
	}

	body := model.Value
	if len(body) > 20 {
		body = body[:20] + "..."
	}

	jobId := strconv.FormatInt(chatModel.JobId, 10)
	chatId := strconv.FormatInt(chatModel.Id, 10)
	pushMsg := universal.PushMessage{
		Title: "New Message",
		Body:  body,
		Link:  "https://naborly.no/chat/" + jobId + "/" + chatId,
	}

	a.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
