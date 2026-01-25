package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

func TestPgGlobalJobs_ByQuery(t *testing.T) {
	ctx := user.WithUser(t.Context(), JobUser)
	teardownSuite, db := SetupJobTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{db: db}
	price := &universal.PriceModel{Value: 10, Currency: "USD"}
	hitDescription, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "I have work for you", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}
	hitCity, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Job in a city", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     &universal.AddressModel{City: "Wojaszowka", PostalCode: "38-471", District: "Podkarpackie", Line1: "Wojaszowka 209", Line2: ""},
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}

	globalJobs := NewPgGlobalJobs(db)
	query := "work"
	foundJobs, err := globalJobs.ByQuery(ctx, &query)
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != hitDescription.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), hitDescription.ID())
	}

	q := "Wojaszowka"
	foundJobs, err = globalJobs.ByQuery(ctx, &q)
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != hitCity.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), hitDescription.ID())
	}
}

func TestPgGlobalJobs_NearBy(t *testing.T) {
	ctx := user.WithUser(t.Context(), JobUser)
	teardownSuite, db := SetupJobTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{db: db}
	price := &universal.PriceModel{Value: 10, Currency: "USD"}

	positionModel := &universal.PositionModel{Lat: 21.6711, Lon: 49.7773}
	job1, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Cut the grass", ImageUrl: "imageUrl"},
			Position:    positionModel,
			Address:     universal.EmptyAddressModel(),
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}
	err = job1.State().Change(ctx, job.JobPublished)
	if err != nil {
		t.Error(err)
	}
	job2, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Job in a country side", ImageUrl: "imageUrl"},
			Position:    positionModel,
			Address:     &universal.AddressModel{City: "Wojaszowka", PostalCode: "38-471", District: "Podkarpackie", Line1: "Wojaszowka 209", Line2: ""},
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}
	err = job2.State().Change(ctx, job.JobPublished)
	if err != nil {
		t.Error(err)
	}

	globalJobs := NewPgGlobalJobs(db)
	foundJobs, err := globalJobs.NearBy(ctx, &universal.RadarModel{Perimeter: 10000, Position: &universal.PositionModel{Lat: 21.6, Lon: 49.7}})
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != job1.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), job1.ID())
	}
}

func TestPgGlobalJobs_Search(t *testing.T) {
	ctx := user.WithUser(t.Context(), JobUser)
	teardownSuite, db := SetupJobTest(t)
	defer teardownSuite(t)

	jobs := PgJobs{db: db}
	price := &universal.PriceModel{Value: 10, Currency: "USD"}

	positionModel := &universal.PositionModel{Lat: 21.6711, Lon: 49.7773}
	job1, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Cut the grass", ImageUrl: "imageUrl"},
			Position:    positionModel,
			Address:     universal.EmptyAddressModel(),
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}
	err = job1.State().Change(ctx, job.JobPublished)
	if err != nil {
		t.Error(err)
	}
	job2, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Remove cutted grass", ImageUrl: "imageUrl"},
			Position:    positionModel,
			Address:     &universal.AddressModel{City: "Wojaszowka", PostalCode: "38-471", District: "Podkarpackie", Line1: "Wojaszowka 209", Line2: ""},
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}
	err = job2.State().Change(ctx, job.JobPublished)
	if err != nil {
		t.Error(err)
	}

	globalJobs := NewPgGlobalJobs(db)

	query := "grass"
	jsi := &job.JobSearchInput{
		Query: &query,
		Radar: &universal.RadarModel{Perimeter: 10000, Position: &universal.PositionModel{Lat: 21.6, Lon: 49.7}},
	}
	foundJobs, err := globalJobs.Search(ctx, jsi)
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != job1.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), job1.ID())
	}
}
