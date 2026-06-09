package job

import (
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type JobInbox interface {
	Messages() universal.Inbox[messages.OLDMessage]
	Offers() universal.Inbox[Offer]
}

type JobOutbox interface {
	Offers() universal.Outbox[Offer]
}
