package user

// API

type Account interface {
	Model() *AccountModel
	Update(newModel *AccountModel) error
}

// Builder

// Model

type AccountModel struct {
	Token string `json:"token"`
}

func (model *AccountModel) Change(newModel *AccountModel) {
	model.Token = newModel.Token
}

// Solid

type SolidAccount struct {
	model   *AccountModel
	Account Account
}

func NewSolidAccount(model *AccountModel, Account Account) Account {
	return &SolidAccount{model, Account}
}

func (addr SolidAccount) Update(newModel *AccountModel) error {
	addr.model.Change(newModel)
	if addr.Account == nil {
		return nil
	}
	return addr.Account.Update(newModel)
}

func (addr SolidAccount) Model() *AccountModel {
	return addr.model
}
