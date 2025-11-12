package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

func TestPgOffers_TestFlow(t *testing.T) {
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

	ctx = user.WithUser(ctx, WorkUser)
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

	waiting, err := newJob.Offers().Waiting()
	if err != nil {
		t.Error(err)
	}
	onlyOffer := waiting[0]
	accepted, err := onlyOffer.Accepted()
	if err != nil {
		t.Error(err)
	}
	if accepted {
		t.Error("accepted")
	}
	rejected, err := onlyOffer.Rejected()
	if err != nil {
		t.Error(err)
	}
	if rejected {
		t.Error("rejected")
	}

	ctx = user.WithUser(ctx, WorkUser)
	err = offer.Reject(ctx)
	if err != nil {
		t.Error(err)
	}

	err = offer.Accept(ctx)
	if err != nil {
		t.Error(err)
	}

}
