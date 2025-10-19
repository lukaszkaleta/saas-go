package universal

// API

type Contents interface {
	ById(id int64) (Content, error)
	ByName(id string) (Content, error)
	Add(model *ContentModel) (Content, error)
}

// Builders

func NewSolidContents(ratings Contents) SolidContents {
	return SolidContents{ratings: ratings}
}

// Solid

type SolidContents struct {
	ratings Contents
}

func (s SolidContents) Add(r ContentModel) (*Content, error) {
	return nil, nil
}

func (s SolidContents) ById(id int64) (Content, error) {
	return s.ratings.ById(id)
}
func (s SolidContents) ByName(name string) (Content, error) {
	return s.ratings.ByName(name)
}
