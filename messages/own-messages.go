package messages

import (
	"context"
)

type Own interface {
	LastQuestionsToMe(ctx context.Context) ([]Message, error)
}
