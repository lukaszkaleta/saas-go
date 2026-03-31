package universal

import "context"

// API
type Person interface {
	Idable
	Model(ctx context.Context) *PersonModel
	Update(ctx context.Context, person *PersonModel) error
	UpdateAverageRating(ctx context.Context, score int) error
	Avatar(ctx context.Context) Description
}

// Builders

// Model

type PersonModel struct {
	Id            int64             `json:"id"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Email         string            `json:"email"`
	Phone         string            `json:"phone"`
	Avatar        *DescriptionModel `json:"avatar"`
	AverageRating int               `json:"rating"`
}

func (model *PersonModel) Change(newModel *PersonModel) {
	model.FirstName = newModel.FirstName
	model.LastName = newModel.LastName
	model.Email = newModel.Email
	model.Phone = newModel.Phone
	model.ChangeAverageRating(newModel.AverageRating)
}

func (model *PersonModel) ChangeAverageRating(score int) {
	model.AverageRating = score
}

func (model *PersonModel) ID() int64 {
	return model.Id
}

func EmptyPersonModel() *PersonModel {
	return &PersonModel{
		FirstName:     "",
		LastName:      "",
		Email:         "",
		Phone:         "",
		AverageRating: 5,
		Avatar:        EmptyDescriptionModel(),
	}
}

func IdEmptyPersonModel(id int64) *PersonModel {
	model := EmptyPersonModel()
	model.Id = id
	return model
}

// Solid

type SolidPerson struct {
	model  *PersonModel
	person Person
}

func NewSolidPerson(model *PersonModel, person Person) Person {
	return &SolidPerson{model, person}
}

func (p SolidPerson) Update(ctx context.Context, newModel *PersonModel) error {
	p.model.Change(newModel)
	if p.person == nil {
		return nil
	}
	return p.person.Update(ctx, newModel)
}

func (p SolidPerson) UpdateAverageRating(ctx context.Context, score int) error {
	p.model.ChangeAverageRating(score)
	if p.person == nil {
		return nil
	}
	return p.person.UpdateAverageRating(ctx, score)
}

func (p SolidPerson) Model(ctx context.Context) *PersonModel {
	return p.model
}

func (p SolidPerson) ID() int64 {
	return p.model.Id
}

func (p SolidPerson) Avatar(ctx context.Context) Description {
	if p.person != nil {
		return NewSolidDescription(p.Model(ctx).Avatar, p.person.Avatar(ctx))
	}
	return NewSolidDescription(p.Model(ctx).Avatar, nil)

}
