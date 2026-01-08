package job

import (
	"context"
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Jobs interface {
	universal.Idables[Job]
	Add(ctx context.Context, model *JobModel) (Job, error)
	List(ctx context.Context) ([]Job, error)
}

func JobModels(jobs []Job) []*JobModel {
	var models []*JobModel
	for _, modelAware := range jobs {
		models = append(models, modelAware.Model()) // note the = instead of :=
	}
	return models
}

func JobHints(jobs []Job) []*JobHint {
	var hints []*JobHint
	for _, o := range jobs {
		if o != nil {
			hints = append(hints, o.Model().Hint()) // note the = instead of :=
		}
	}
	return hints
}

func GeoJobs(jobs []Job) universal.GeoFeatureCollection[JobHint] {
	features := make([]universal.GeoFeature[JobHint], 0, len(jobs))
	for i := range jobs {
		m := jobs[i]
		pt := universal.NewGeoPoint(m.Model().Position.Lon, m.Model().Position.Lat)
		features = append(features, universal.NewGeoFeature[JobHint](strconv.FormatInt(m.Model().Id, 10), pt, *m.Model().Hint()))
	}
	return universal.NewGeoFeatureCollection(features)
}
