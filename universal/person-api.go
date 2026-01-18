package universal

import "context"

// API
type Person interface {
	Idable
	Model(ctx context.Context) *PersonModel
	Update(ctx context.Context, person *PersonModel) error
	Avatar(ctx context.Context) Description
	Ratings() Ratings
}

// Builders

// Model

type PersonModel struct {
	Id            int64             `json:"id"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Email         string            `json:"email"`
	Phone         string            `json:"phone"`
	Avatar        *DescriptionModel `json:"avatarUrl"`
	AverageRating int               `json:"rating"`
}

func (model *PersonModel) Change(newModel *PersonModel) {
	model.FirstName = newModel.FirstName
	model.LastName = newModel.LastName
	model.Email = newModel.Email
	model.Phone = newModel.Phone
}

func EmptyPersonModel() *PersonModel {
	return &PersonModel{
		FirstName:     "",
		LastName:      "",
		Email:         "",
		Phone:         "",
		AverageRating: 10,
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

func (p SolidPerson) Model(ctx context.Context) *PersonModel {
	return p.model
}

func (p SolidPerson) ID() int64 {
	return p.model.Id
}

func (p SolidPerson) Ratings() Ratings {
	if p.person == nil {
		return DummyRatings{}
	}
	return p.person.Ratings()
}

func (p SolidPerson) Avatar(ctx context.Context) Description {
	if p.person != nil {
		return NewSolidDescription(p.Model(ctx).Avatar, p.person.Avatar(ctx))
	}
	return NewSolidDescription(p.Model(ctx).Avatar, nil)

}
