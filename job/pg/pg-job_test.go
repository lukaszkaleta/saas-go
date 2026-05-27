package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
	pguser "github.com/lukaszkaleta/saas-go/user/pg"
)

var JobUser = user.WithId(1)
var WorkUser = user.WithId(2)

func SetupJobTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "job_test")
	// Ensure filestore_filesystem exists (or at least avoid the FK error)
	_, err := db.Pool.Exec(tb.Context(), "CREATE TABLE if not exists filestore_filesystem (id serial primary key, name_value text not null default '', name_slug text not null default '')")
	if err != nil {
		tb.Fatal(err)
	}

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

	// Ensure filestore_filesystem exists (or at least avoid the FK error)
	_, err = db.Pool.Exec(tb.Context(), "CREATE TABLE if not exists filestore_filesystem (id serial primary key, name_value text not null default '', name_slug text not null default '')")
	if err != nil {
		tb.Fatal(err)
	}

	err = fsSchema.CreateTest()
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
			PriceOwner:  &job.PriceFormula{Mode: job.FIXED, Value: 100},
			PriceWorker: &job.PriceFormula{Mode: job.PERCENT, Value: 5},
			Tags:        []string{"tag1", "tag2"},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
	globalJobs := NewPgGlobalJobs(db, pguser.NewPgUsers(db).Search())
	model, err := newJob.Model(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if model.PriceOwner.Mode != job.FIXED || model.PriceOwner.Value != 100 {
		t.Errorf("PriceOwner mismatch: got %+v, want Mode: FIXED, Value: 100", model.PriceOwner)
	}
	if model.PriceWorker.Mode != job.PERCENT || model.PriceWorker.Value != 5 {
		t.Errorf("PriceWorker mismatch: got %+v, want Mode: PERCENT, Value: 5", model.PriceWorker)
	}

	jobById, err := globalJobs.ById(t.Context(), model.Id)
	if err != nil {
		t.Error(err)
	}
	if jobById == nil {
		t.Error("job by id should not be nil")
	} else {
		modelById, _ := jobById.Model(t.Context())
		if modelById.PriceOwner.Mode != job.FIXED || modelById.PriceOwner.Value != 100 {
			t.Errorf("PriceOwner mismatch in jobById: got %+v", modelById.PriceOwner)
		}
		if modelById.PriceWorker.Mode != job.PERCENT || modelById.PriceWorker.Value != 5 {
			t.Errorf("PriceWorker mismatch in jobById: got %+v", modelById.PriceWorker)
		}
	}
}
