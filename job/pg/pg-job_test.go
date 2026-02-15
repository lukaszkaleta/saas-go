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

func SetupJobTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "job_test")
	fsSchema := pgfilestore.NewFilestoreSchema(db)
	schema := NewJobSchema(db)

	dropFunc := func(tb testing.TB) {
		err := schema.DropTest()
		if err != nil {
			panic(err)
		}
		err = fsSchema.DropTest()
		if err != nil {
			panic(err)
		}
	}

	dropFunc(tb)

	err := fsSchema.CreateTest()
	if err != nil {
		tb.Fatal(err)
	}
	err = schema.CreateTest()
	if err != nil {
		tb.Fatal(err)
	}

	_, err = db.Pool.Exec(tb.Context(), "delete from job")
	if err != nil {
		tb.Fatal(err)
	}
	_, err = db.Pool.Exec(tb.Context(), "delete from users")
	if err != nil {
		tb.Error(err)
	}
	for i := 1; i < 3; i++ {
		_, err := db.Pool.Exec(tb.Context(), "insert into users (id) values ($1)", i)
		if err != nil {
			tb.Error(err)
		}
	}

	return dropFunc, db
}

func TestPgJob_Status(t *testing.T) {
	ctx := user.WithUser(t.Context(), JobUser)
	teardownSuite, db := SetupJobTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{db: db}
	newJob, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "description", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       universal.EmptyPriceModel(),
			Tags:        []string{"tag1", "tag2"},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
	if job.JobDraft != newJob.State().Name(t.Context()) {
		t.Error("job status is not draft")
	}
	globalJobs := NewPgGlobalJobs(db)
	model, err := newJob.Model(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	jobById, err := globalJobs.ById(t.Context(), model.Id)
	if err != nil {
		t.Error(err)
	}
	if jobById == nil {
		t.Error("job by id should not be nil")
	}
	if job.JobDraft != jobById.State().Name(t.Context()) {
		t.Error("job status is not draft")
	}
}
