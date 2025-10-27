package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

func setupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas", "job-test")
	fsSchema := pgfilestore.NewFilestoreSchema(db)
	fsSchema.Create()
	schema := NewJobSchema(db)
	schema.Create()

	return func(tb testing.TB) {
		schema.Drop()
		fsSchema.Drop()
	}, db
}

func TestPgJob_Status(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{Db: db}
	newJob, err := jobs.Add(
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "description", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       &universal.PriceModel{},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if job.JobDraft != newJob.State().Name() {
		t.Error("job status is not draft")
	}
	globalJobs := PgGlobalJobs{Db: db}
	jobById, err := globalJobs.ById(newJob.Model().Id)
	if err != nil {
		t.Error(err)
	}
	if job.JobDraft != jobById.State().Name() {
		t.Error("job status is not draft")
	}

}
