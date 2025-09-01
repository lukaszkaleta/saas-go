package pgoffer

import (
	"context"

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
	query := "select * from offer"
	offers := []offer.Offer{}
	rows, err := globalOffers.Db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		descriptionRow := new(universal.DescriptionModel)
		addressRow := new(universal.AddressModel)
		priceRow := new(universal.PriceModel)
		positionRow := new(universal.PositionModel)
		err := rows.Scan(
			&id,
			&descriptionRow.Value,
			&descriptionRow.ImageUrl,
			&addressRow.Line1,
			&addressRow.Line2,
			&addressRow.City,
			&addressRow.PostalCode,
			&addressRow.District,
			&positionRow.Lat,
			&positionRow.Lon,
			&priceRow.Value,
			&priceRow.Currency,
		)
		pgOffer := &PgOffer{Db: globalOffers.Db, Id: id}
		if err != nil {
			return nil, err
		}

		solidOffer := offer.NewSolidOffer(
			&offer.OfferModel{
				Id:          id,
				Description: descriptionRow,
				Address:     addressRow,
				Price:       priceRow,
				Position:    positionRow,
			},
			pgOffer,
			id)
		offers = append(offers, solidOffer)
	}
	return offers, nil
}

func (globalOffers *PgGlobalOffers) ById(id int64) (offer.Offer, error) {
	query := "select * from offer where id = $1"
	rows, err := globalOffers.Db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return nil, nil
	}

	descriptionRow := new(universal.DescriptionModel)
	addressRow := new(universal.AddressModel)
	priceRow := new(universal.PriceModel)
	positionRow := new(universal.PositionModel)
	err = rows.Scan(
		&id,
		&descriptionRow.Value,
		&descriptionRow.ImageUrl,
		&addressRow.Line1,
		&addressRow.Line2,
		&addressRow.City,
		&addressRow.PostalCode,
		&addressRow.District,
		&positionRow.Lat,
		&positionRow.Lon,
		&priceRow.Value,
		&priceRow.Currency,
	)
	pgOffer := &PgOffer{Db: globalOffers.Db, Id: id}
	if err != nil {
		return nil, err
	}

	solidOffer := offer.NewSolidOffer(
		&offer.OfferModel{
			Id:          id,
			Description: descriptionRow,
			Address:     addressRow,
			Price:       priceRow,
			Position:    positionRow,
		},
		pgOffer,
		id)

	return solidOffer, nil
}
