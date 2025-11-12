package user

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type User interface {
	universal.Idable
	Model() *UserModel
	Account() Account
	Person() universal.Person
	Address() universal.Address
	Settings() UserSettings
	FileSystem(name string) (filestore.FileSystem, error)
	Archive() error
}

func WithUser(ctx context.Context, usr User) context.Context {
	return context.WithValue(ctx, "user-context", usr.Model())
}

func FetchUser(ctx context.Context) *UserModel {
	return ctx.Value("user-context").(*UserModel)
}

func WithId(id int64) User {
	model := NewUserModel()
	model.Id = id
	return NewSolidUser(model, nil)
}

type UserFs string

const (
	UserAvatarFs UserFs = "user-avatar"
)

// Model

type UserModel struct {
	Id       int64                   `json:"id"`
	Person   *universal.PersonModel  `json:"person"`
	Address  *universal.AddressModel `json:"address"`
	Account  *AccountModel           `json:"account"`
	Settings *UserSettingsModel      `json:"settings"`
}

func NewUserModel() *UserModel {
	return &UserModel{
		Id:       0,
		Person:   &universal.PersonModel{},
		Address:  &universal.AddressModel{},
		Account:  &AccountModel{},
		Settings: NewUserSettingsModel(),
	}
}

func EmptyUserModel() *UserModel {
	return &UserModel{
		Id:       0,
		Person:   universal.EmptyPersonModel(),
		Address:  universal.EmptyAddressModel(),
		Account:  EmptyAccountModel(),
		Settings: EmptyUserSettingsModel(),
	}
}

// Solid

type SolidUser struct {
	Id    int64
	model *UserModel
	user  User
}

func NewSolidUser(model *UserModel, user User) User {
	return &SolidUser{
		model.Id,
		model,
		user,
	}
}

func (u SolidUser) Model() *UserModel {
	return u.model
}

func (u SolidUser) Person() universal.Person {
	if u.user != nil {
		return universal.NewSolidPerson(
			u.Model().Person,
			u.user.Person(),
		)
	}
	return universal.NewSolidPerson(u.Model().Person, nil)
}

func (u SolidUser) Address() universal.Address {
	if u.user != nil {
		return universal.NewSolidAddress(
			u.Model().Address,
			u.user.Address(),
		)
	}
	return universal.NewSolidAddress(u.Model().Address, nil)
}

func (u SolidUser) Settings() UserSettings {
	return u.user.Settings()
}

func (u SolidUser) Account() Account {
	return u.user.Account()
}

func (u SolidUser) FileSystem(name string) (filestore.FileSystem, error) {
	return u.user.FileSystem(name)
}

func (u SolidUser) Archive() error {
	return nil
}

func (u SolidUser) ID() int64 {
	return u.Id
}

// Relations
