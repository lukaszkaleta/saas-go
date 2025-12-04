package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

var JobUser = user.WithId(1)
var WorkUser = user.WithId(2)

func setupJobTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "job-test")
	fsSchema := pgfilestore.NewFilestoreSchema(db)
	err := fsSchema.CreateTest()
	if err != nil {
		tb.Fatal(err)
	}
	schema := NewJobSchema(db)
	err = schema.CreateTest()
	if err != nil {
		tb.Fatal(err)
	}

	for i := 1; i < 3; i++ {
		_, err := db.Pool.Exec(tb.Context(), "insert into users (id) values ($1)", i)
		if err != nil {
			tb.Error(err)
		}
	}

	return func(tb testing.TB) {
		err := schema.DropTest()
		if err != nil {
			panic(err)
		}
		err = fsSchema.DropTest()
		if err != nil {
			panic(err)
		}
	}, db
}

func TestPgJob_Status(t *testing.T) {
	ctx := user.WithUser(t.Context(), JobUser)
	teardownSuite, db := setupJobTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{Db: db}
	newJob, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "description", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       universal.EmptyPriceModel(),
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
	if jobById == nil {
		t.Error("job by id should not be nil")
	}
	if job.JobDraft != jobById.State().Name() {
		t.Error("job status is not draft")
	}

}
