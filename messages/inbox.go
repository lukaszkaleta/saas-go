package messages

import (
	"context"
)

type Inbox interface {
	LastQuestions(ctx context.Context) ([]Message, error)
	LastAnswers(ctx context.Context) ([]Message, error)
}
