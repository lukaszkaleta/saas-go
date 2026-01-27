package user

import "context"

type UserSearch interface {
	ByPhone(ctx context.Context, phone string) (User, error)
	ModelsByIds(ctx context.Context, ids []int64) ([]*UserModel, error)
}
