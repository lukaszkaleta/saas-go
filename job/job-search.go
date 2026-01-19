package job

import "github.com/lukaszkaleta/saas-go/universal"

type JobSearch interface {
	universal.FullText[Job]
}
