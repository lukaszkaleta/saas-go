package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

func TestPgOffers_Make(t *testing.T) {
	teardownSuite, db := setupJobTest(t)
	defer teardownSuite(t)
	ctx := user.WithUser(t.Context(), JobUser)

	jobs := PgJobs{Db: db}
	price := &universal.PriceModel{Value: 10, Currency: "USD"}
	newJob, err := jobs.Add(
		ctx,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "description", ImageUrl: "imageUrl"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       price,
		},
	)
	if err != nil {
		t.Error(err)
	}

	offerModel := &job.OfferModel{
		Description: &universal.DescriptionModel{Value: "I will do it"},
		Price:       &universal.PriceModel{Value: price.Value - 1, Currency: price.Currency},
		Rating:      10,
	}
	offer, err := newJob.Offers().Make(ctx, offerModel)
	if err != nil {
		t.Error(err)
	}
	if offer == nil {
		t.Error("offer was nil")
	}
}
