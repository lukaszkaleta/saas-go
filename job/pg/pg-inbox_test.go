package pgjob

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

func TestPgOfferInbox_Last(t *testing.T) {
	teardownSuite, db := SetupJobTest(t)
	defer teardownSuite(t)

	// User 1 creates a job
	ctx1 := user.WithUser(t.Context(), JobUser)
	jobs := PgJobs{db: db}
	price := &universal.PriceModel{Value: 100, Currency: "USD"}
	newJob, err := jobs.Add(
		ctx1,
		&job.JobModel{
			Description: &universal.DescriptionModel{Value: "Job by User 1"},
			Position:    &universal.PositionModel{Lon: 1, Lat: 1},
			Address:     universal.EmptyAddressModel(),
			Price:       price,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// User 2 makes an offer for User 1's job
	ctx2 := user.WithUser(t.Context(), WorkUser)
	offerModel := &job.OfferModel{
		Description: &universal.DescriptionModel{Value: "Offer by User 2"},
		Price:       &universal.PriceModel{Value: 90, Currency: "USD"},
	}
	_, err = newJob.Offers().Make(ctx2, offerModel)
	if err != nil {
		t.Fatal(err)
	}

	// User 1 checks their inbox
	inbox := NewPgOfferInbox(db)

	// Count unread
	count, err := inbox.CountUnread(ctx1)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("expected 1 unread offer, got %d", count)
	}

	// Last offers
	offers, err := inbox.Last(ctx1)
	if err != nil {
		t.Error(err)
	}
	if len(offers) != 1 {
		t.Errorf("expected 1 offer in inbox, got %d", len(offers))
	}

	// User 2 checks their inbox (should be empty as they don't own the job)
	count2, err := inbox.CountUnread(ctx2)
	if err != nil {
		t.Error(err)
	}
	if count2 != 0 {
		t.Errorf("expected 0 unread offers for User 2, got %d", count2)
	}
}
