package pgoffer

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/offer"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgGlobalOffers struct {
	Db *pg.PgDb
}

func NewPgGlobalOffers(Db *pg.PgDb) offer.GlobalOffers {
	return &PgGlobalOffers{Db}
}

func (globalOffers *PgGlobalOffers) NearBy(radar *universal.RadarModel) ([]offer.Offer, error) {
	id := int64(0)
	query := "select * from job"
	offers := []offer.Offer{}
	rows, err := globalOffers.Db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pgOffer := &PgOffer{Db: globalOffers.Db, Id: id}
		offerModel, err := MapOffer(rows)
		if err != nil {
			return nil, err
		}
		solidOffer := offer.NewSolidOffer(
			offerModel,
			pgOffer,
			id)
		offers = append(offers, solidOffer)
	}
	return offers, nil
}

func (globalOffers *PgGlobalOffers) ById(id int64) (offer.Offer, error) {
	query := "select * from job where id = @id"
	rows, err := globalOffers.Db.Pool.Query(context.Background(), query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}

	pgOffer := &PgOffer{Db: globalOffers.Db, Id: id}
	offerModel, err := MapOffer(rows)
	if err != nil {
		return nil, err
	}

	return offer.NewSolidOffer(offerModel, pgOffer, offerModel.Id), nil
}
