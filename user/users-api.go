package user

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Users interface {
	Add(context context.Context, model *universal.PersonModel) (User, error)
	ById(context context.Context, id int64) (User, error)
	ListAll(context context.Context) ([]User, error)
	Search() UserSearch
	EstablishAccount(context context.Context, model *UserModel) (User, error)
}

func UserModels(ctx context.Context, users []User) ([]*UserModel, error) {
	var models []*UserModel
	for _, u := range users {
		model, err := u.Model(ctx)
		if err != nil {
			return nil, err
		}
		models = append(models, model) // note the = instead of :=
	}
	return models, nil
}
