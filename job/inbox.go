package job

import (
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

type JobInbox interface {
	Messages() universal.Inbox[messages.Message]
	Offers() universal.Inbox[Offer]
}
