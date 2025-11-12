package universal

// API
type Person interface {
	ID() int64
	Model() *PersonModel
	Update(person *PersonModel) error

	Ratings() Ratings
}

// Builders

// Model

type PersonModel struct {
	Id            int64  `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	AverageRating int    `json:"rating"`
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
	}
}

// Solid

type SolidPerson struct {
	model  *PersonModel
	person Person
}

func NewSolidPerson(model *PersonModel, person Person) Person {
	return &SolidPerson{model, person}
}

func (p SolidPerson) Update(newModel *PersonModel) error {
	p.model.Change(newModel)
	if p.person == nil {
		return nil
	}
	return p.person.Update(newModel)
}

func (p SolidPerson) Model() *PersonModel {
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
