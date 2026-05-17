package user

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type User interface {
	universal.Idable
	Model(ctx context.Context) (*UserModel, error)
	Account() Account
	Person() universal.Person
	Address() universal.Address
	Settings() UserSettings
	FileSystem(name string) (filestore.FileSystem, error)
	Rated() universal.Rated
	Archive() error
}

func WithUser(ctx context.Context, usr User) context.Context {
	model, err := usr.Model(ctx)
	if err != nil {
		return ctx
	}
	ctx = context.WithValue(ctx, "current-user", model)
	id := usr.ID()
	return context.WithValue(ctx, universal.CurrentUserKey, &id)
}

func CurrentUser(ctx context.Context) *UserModel {
	return ctx.Value("current-user").(*UserModel)
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

func (u UserModel) ID() int64 {
	return u.Id
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

func (u SolidUser) Model(ctx context.Context) (*UserModel, error) {
	return u.model, nil
}

func (u SolidUser) Person() universal.Person {
	if u.user != nil {
		model, err := u.Model(context.Background())
		if err != nil {
			panic(err)
		}
		return universal.NewSolidPerson(
			model.Person,
			u.user.Person(),
		)
	}
	model, err := u.Model(context.Background())
	if err != nil {
		panic(err)
	}
	return universal.NewSolidPerson(model.Person, nil)
}

func (u SolidUser) Address() universal.Address {
	if u.user != nil {
		model, err := u.Model(context.Background())
		if err != nil {
			panic(err)
		}
		return universal.NewSolidAddress(
			model.Address,
			u.user.Address(),
		)
	}
	model, err := u.Model(context.Background())
	if err != nil {
		panic(err)
	}
	return universal.NewSolidAddress(model.Address, nil)
}

func (u SolidUser) Settings() UserSettings {
	return u.user.Settings()
}

func (u SolidUser) Account() Account {
	return u.user.Account()
}

func (u SolidUser) Rated() universal.Rated {
	if u.user != nil {
		return u.user.Rated()
	}
	return universal.DummyRatings{}
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
