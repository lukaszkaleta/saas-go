package job

import "github.com/lukaszkaleta/saas-go/universal"

type JobSearchPaging struct {
	Id   int64   `json:"id"`
	Rank float64 `json:"rank"`
}

type JobSearchInput struct {
	// query from client
	Query *string `json:"query"`
	// radar from client
	Radar *universal.RadarModel `json:"radar"`
	// paging from server (send back and forth from client)
	Paging JobSearchPaging `json:"paging"`
}

type JobSearchRanking struct {
	Distance *float64 `json:"distance"`
	Rank     *float64 `json:"rank"`
}

type JobSearchOutput struct {
	Model   *JobModel         `json:"job"`
	Ranking *JobSearchRanking `json:"ranking"`
	Paging  *JobSearchPaging  `json:"paging"`
}

func (jobSearchOutput JobSearchOutput) ID() int64 {
	return jobSearchOutput.Model.Id
}
