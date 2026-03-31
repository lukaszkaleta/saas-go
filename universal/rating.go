package universal

import "context"

// API
type Rating interface {
	Idable
	ActionsAware
	RevieweeId(ctx context.Context) int64
	SubjectId(ctx context.Context) int64
	Model(ctx context.Context) (*RatingModel, error)
	Update(ctx context.Context, rating *RatingModel) error
}

// Model

type RatingModel struct {
	Id         int64             `json:"id"`
	RevieweeId int64             `json:"revieweeId"`
	SubjectId  int64             `json:"subjectId"`
	Review     *DescriptionModel `json:"review"`
	Score      int               `json:"score"`
	Actions    *ActionsModel     `json:"actions"`
}

func (m *RatingModel) Change(newModel *RatingModel) {
	m.Score = newModel.Score
	m.Review = newModel.Review
}

func EmptyRatingModel() *RatingModel {
	return &RatingModel{
		Review:     EmptyDescriptionModel(),
		Score:      0,
		RevieweeId: 0,
		SubjectId:  0,
		Actions:    EmptyActionsModel(),
	}
}

// Solid

type SolidRating struct {
	Id     int64
	model  *RatingModel
	rating Rating
}

func (s *SolidRating) ID() int64 {
	return s.Id
}

func (s *SolidRating) Actions() Actions {
	return s.rating.Actions()
}

func (s *SolidRating) RevieweeId(ctx context.Context) int64 {
	return s.model.RevieweeId
}

func (s *SolidRating) SubjectId(ctx context.Context) int64 {
	return s.model.SubjectId
}

func NewSolidRating(ratingModel *RatingModel, rating Rating, id int64) Rating {
	return &SolidRating{
		model:  ratingModel,
		rating: rating,
		Id:     id,
	}
}

func (s *SolidRating) Model(ctx context.Context) (*RatingModel, error) {
	return s.model, nil
}

func (s *SolidRating) Update(ctx context.Context, newModel *RatingModel) error {
	s.model.Change(newModel)
	if s.rating != nil {
		return s.rating.Update(ctx, newModel)
	}
	return nil
}

// Dummy

type DummyRating struct {
}

func (d DummyRating) ID() int64 {
	return 0
}

func (d DummyRating) Actions() Actions {
	return &SolidActions{}
}

func (d DummyRating) RevieweeId(ctx context.Context) int64 {
	return 0
}

func (d DummyRating) SubjectId(ctx context.Context) int64 {
	return 0
}

func (d DummyRating) Model(ctx context.Context) (*RatingModel, error) {
	return EmptyRatingModel(), nil
}

func (d DummyRating) Update(ctx context.Context, rating *RatingModel) error {
	return nil
}
