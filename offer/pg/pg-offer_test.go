package pgoffer

import (
	"log"
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/offer"
	"github.com/lukaszkaleta/saas-go/universal"
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test suite")
	return func(tb testing.TB) {
		log.Println("teardown test suite")
	}
}

func setupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("offer_test")
	fsSchema := pgfilestore.NewFilestoreSchema(db)
	fsSchema.Create()
	schema := NewOfferSchema(db)
	schema.Create()

	return func(tb testing.TB) {
		schema.Drop()
		fsSchema.Drop()
	}, db
}

func TestPgOffer_Status(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	offers := PgOffers{Db: db}
	newOffer, err := offers.AddWithPlace(
		&universal.PositionModel{Lon: 1, Lat: 1},
		universal.EmptyAddressModel(),
	)
	if err != nil {
		t.Error(err)
	}
	if offer.OfferDraft != newOffer.State().Name() {
		t.Error("offer status is not draft")
	}
	globalOffers := PgGlobalOffers{Db: db}
	offerById, err := globalOffers.ById(newOffer.Model().Id)
	if err != nil {
		t.Error(err)
	}
	if offer.OfferDraft != offerById.State().Name() {
		t.Error("offer status is not draft")
	}

}
