package user

import "context"

// API

type Account interface {
	Model(ctx context.Context) (*AccountModel, error)
	Update(ctx context.Context, newModel *AccountModel) error
	UpdatePushNotificationToken(ctx context.Context, token string) error
}

// Builder

// Model

type AccountModel struct {
	Token         string `json:"token"`
	FirebaseToken string `json:"firebaseToken" omitzero:"true"`
}

func (model *AccountModel) Change(newModel *AccountModel) {
	model.Token = newModel.Token
	model.FirebaseToken = newModel.FirebaseToken
}

func EmptyAccountModel() *AccountModel {
	return &AccountModel{
		Token:         "",
		FirebaseToken: "",
	}
}

// Solid

type SolidAccount struct {
	model   *AccountModel
	Account Account
}

func NewSolidAccount(model *AccountModel, Account Account) Account {
	return &SolidAccount{model, Account}
}

func (addr SolidAccount) Update(ctx context.Context, newModel *AccountModel) error {
	addr.model.Change(newModel)
	if addr.Account == nil {
		return nil
	}
	return addr.Account.Update(ctx, newModel)
}

func (addr SolidAccount) UpdatePushNotificationToken(ctx context.Context, token string) error {
	addr.model.FirebaseToken = token
	if addr.Account == nil {
		return nil
	}
	return addr.Account.UpdatePushNotificationToken(ctx, token)
}

func (addr SolidAccount) Model(ctx context.Context) (*AccountModel, error) {
	return addr.model, nil
}
