package messages

import (
	"context"
)

type Inbox interface {
	LastQuestions(ctx context.Context) ([]Message, error)
	CountUnreadQuestions(ctx context.Context) (int, error)
	LastAnswers(ctx context.Context) ([]Message, error)
	CountUnreadAnswers(ctx context.Context) (int, error)
}
