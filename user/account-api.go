package user

import "context"

// API

type Account interface {
	Model(ctx context.Context) *AccountModel
	Update(ctx context.Context, newModel *AccountModel) error
}

// Builder

// Model

type AccountModel struct {
	Token string `json:"token"`
}

func (model *AccountModel) Change(newModel *AccountModel) {
	model.Token = newModel.Token
}

func EmptyAccountModel() *AccountModel {
	return &AccountModel{
		Token: "",
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

func (addr SolidAccount) Model(ctx context.Context) *AccountModel {
	return addr.model
}
