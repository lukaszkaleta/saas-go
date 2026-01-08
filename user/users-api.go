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

func UserModels(users []User) []*UserModel {
	var models []*UserModel
	for _, u := range users {
		models = append(models, u.Model(context.Background())) // note the = instead of :=
	}
	return models
}
