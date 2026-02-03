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

func JobModels(ctx context.Context, jobs []Job) []*JobModel {
	var models []*JobModel
	for _, modelAware := range jobs {
		model, _ := modelAware.Model(ctx)
		models = append(models, model) // note the = instead of :=
	}
	return models
}

func JobHints(ctx context.Context, jobs []Job) []*JobHint {
	var hints []*JobHint
	for _, o := range jobs {
		if o != nil {
			model, _ := o.Model(ctx)
			hints = append(hints, model.Hint()) // note the = instead of :=
		}
	}
	return hints
}

func GeoJobs(jobs []*JobSearchOutput) universal.GeoFeatureCollection[JobHint] {
	features := make([]universal.GeoFeature[JobHint], 0, len(jobs))
	for i := range jobs {
		m := jobs[i]
		pt := universal.NewGeoPoint(m.Model.Position.Lon, m.Model.Position.Lat)
		features = append(features, universal.NewGeoFeature[JobHint](strconv.FormatInt(m.Model.Id, 10), pt, *m.Model.Hint()))
	}
	return universal.NewGeoFeatureCollection(features)
}
