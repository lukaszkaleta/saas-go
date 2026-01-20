package job

import "github.com/lukaszkaleta/saas-go/universal"

type JobSearch interface {
}

type JobSearchPosition struct {
	Id   int64   `db:"id"`
	Rank float64 `json:"rank"`
}

type JobSearchInput struct {
	Query    string                `json:"query"`
	Radar    *universal.RadarModel `json:"radar"`
	Position JobSearchPosition     `json:"position"`
}

type JobSearchResult struct {
	Distance int     `json:"distance"`
	Rank     float64 `json:"rank"`
}

type JobSearchModel struct {
	Model     *JobModel       `json:"model"`
	JobSearch JobSearchResult `json:"search"`
}
