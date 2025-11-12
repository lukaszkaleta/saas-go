package universal

// API

type Ratings interface {
	ById(id int64) (Rating, error)
	Add(model *RatingModel) (Rating, error)
}

// Builders

func NewSolidRatings(ratings Ratings) SolidRatings {
	return SolidRatings{ratings: ratings}
}

// Solid

type SolidRatings struct {
	ratings Ratings
}

func (s SolidRatings) Add(r RatingModel) (*Rating, error) {
	return nil, nil
}

func (s SolidRatings) ById(id int64) (Rating, error) {
	return s.ratings.ById(id)
}

// Dummy

type DummyRatings struct {
}

func (dummy DummyRatings) Add(r *RatingModel) (Rating, error) {
	return nil, nil
}

func (dummy DummyRatings) ById(id int64) (Rating, error) {
	return DummyRating{}, nil
}
