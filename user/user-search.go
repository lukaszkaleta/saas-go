package user

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type UserSearch interface {
	ByPhone(ctx context.Context, phone string) (User, error)
	PersonModelsByIds(ctx context.Context, ids []int64) ([]*universal.PersonModel, error)
	ModelsByIds(ctx context.Context, ids []int64) ([]*UserModel, error)
}
