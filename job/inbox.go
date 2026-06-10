package job

import (
	"github.com/lukaszkaleta/saas-go/chat"
	"github.com/lukaszkaleta/saas-go/universal"
)

type JobInbox interface {
	Messages() universal.Inbox[chat.Message]
	Offers() universal.Inbox[Offer]
}

type JobOutbox interface {
	Offers() universal.Outbox[Offer]
}
