package universal

import "context"

// API

type Rated interface {
	Average(ctx context.Context) (int, error)
	AllModels(ctx context.Context) ([]*RatingModel, error)
}

type Ratings interface {
	Adder[*RatingModel, Rating]
	ById(ctx context.Context, id int64) (Rating, error)
	Rated
}

// Builders

func NewSolidRatings(ratings Ratings) Ratings {
	return &SolidRatings{ratings: ratings}
}

// Solid

type SolidRatings struct {
	ratings Ratings
}

func (s SolidRatings) Add(ctx context.Context, r *RatingModel) (Rating, error) {
	return s.ratings.Add(ctx, r)
}

func (s SolidRatings) ById(ctx context.Context, id int64) (Rating, error) {
	return s.ratings.ById(ctx, id)
}

func (s SolidRatings) Average(ctx context.Context) (int, error) {
	if s.ratings == nil {
		return 0, nil
	}
	return s.ratings.Average(ctx)
}

func (s SolidRatings) AllModels(ctx context.Context) ([]*RatingModel, error) {
	if s.ratings == nil {
		return nil, nil
	}
	return s.ratings.AllModels(ctx)
}

// Dummy

type DummyRatings struct {
}

func (dummy DummyRatings) Add(ctx context.Context, r *RatingModel) (Rating, error) {
	return nil, nil
}

func (dummy DummyRatings) ById(ctx context.Context, id int64) (Rating, error) {
	return DummyRating{}, nil
}

func (dummy DummyRatings) Average(ctx context.Context) (int, error) {
	return 0, nil
}

func (dummy DummyRatings) AllModels(ctx context.Context) ([]*RatingModel, error) {
	return nil, nil
}
