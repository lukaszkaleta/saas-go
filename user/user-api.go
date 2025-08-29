package user

import (
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type User interface {
	Model() *UserModel
	Account() Account
	Person() universal.Person
	Address() universal.Address
	Settings() UserSettings
	FileSystem(name string) (filestore.FileSystem, error)
	Archive() error
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

// Solid

type SolidUser struct {
	Id    int64
	model *UserModel
	user  User
}

func NewSolidUser(model *UserModel, user User, id int64) User {
	return &SolidUser{
		id,
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

// Relations
