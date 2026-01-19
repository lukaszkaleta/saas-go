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
	foundJobs, err := globalJobs.ByQuery(ctx, "work")
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != hitDescription.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), hitDescription.ID())
	}

	foundJobs, err = globalJobs.ByQuery(ctx, "Wojaszowka")
	if err != nil {
		t.Fatal(err)
	}
	if foundJobs[0].ID() != hitCity.ID() {
		t.Errorf("got %v, want %v", foundJobs[0].ID(), hitDescription.ID())
	}
}
