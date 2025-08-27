package user

import (
	"github.com/lukaszkaleta/saas-go/universal"
)

type Users interface {
	Add(model *universal.PersonModel) (User, error)
	ById(id int64) (User, error)
	ListAll() ([]User, error)
	Search() UserSearch
	EstablishAccount(model *UserModel) (User, error)
}

func UserModels(users []User) []*UserModel {
	var models []*UserModel
	for _, u := range users {
		models = append(models, u.Model()) // note the = instead of :=
	}
	return models
}
