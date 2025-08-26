package universal

// API
type Rating interface {
	Model() *RatingModel
	Update(rating *RatingModel) error
}

// Model

type RatingModel struct {
	Description string
	Score       int
}

func (m *RatingModel) Change(newModel *RatingModel) {
	m.Score = newModel.Score
	m.Description = newModel.Description
}

// Solid

type SolidRating struct {
	Id     int64
	model  *RatingModel
	rating Rating
}

func NewSolidRating(ratingModel *RatingModel, rating Rating, id int64) Rating {
	return &SolidRating{
		model:  ratingModel,
		rating: rating,
		Id:     id,
	}
}

func (s *SolidRating) Model() *RatingModel {
	return s.model
}

func (s *SolidRating) Update(newModel *RatingModel) error {
	s.model.Change(newModel)
	return nil
}
