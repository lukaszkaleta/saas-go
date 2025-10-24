package pgoffer

import (
	"context"
	"fmt"
	"time"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/offer"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgOffers struct {
	Db  *pg.PgDb
	Ids []int
}

func (pgOffers *PgOffers) Add(model *offer.OfferModel) (offer.Offer, error) {
	offerId := int64(0)
	query := "INSERT INTO job (description_value, description_image_url, position_latitude, position_longitude, address_line_1, address_line_2, address_city, address_postal_code, address_district, price_value, price_currency) VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 ) returning id"
	row := pgOffers.Db.Pool.QueryRow(
		context.Background(),
		query,
		model.Description.Value,
		model.Description.ImageUrl,
		model.Position.Lat,
		model.Position.Lon,
		model.Address.Line1,
		model.Address.Line2,
		model.Address.City,
		model.Address.PostalCode,
		model.Address.District,
		model.Price.Value,
		model.Price.Currency,
	)
	err := row.Scan(&offerId)
	if err != nil {
		return nil, err
	}
	pgOffer := PgOffer{
		Db: pgOffers.Db,
		Id: offerId,
	}
	return offer.NewSolidOffer(
		&offer.OfferModel{
			Id:          offerId,
			Description: &universal.DescriptionModel{},
			Position:    model.Position,
			Address:     model.Address,
			Price:       &universal.PriceModel{},
			State:       offer.OfferStatus{Draft: time.Now()},
		},
		&pgOffer,
		offerId,
	), nil
}

// Relation

type PgRelationJobs struct {
	Db       *pg.PgDb
	Offers   *PgOffers
	Relation pg.RelationEntity
}

func NewPgRelationJobs(pfOffers *PgOffers, relation pg.RelationEntity) PgRelationJobs {
	return PgRelationJobs{
		Db:       pfOffers.Db,
		Offers:   pfOffers,
		Relation: relation,
	}
}

func (p PgRelationJobs) Add(offerModel *offer.OfferModel) (offer.Offer, error) {
	newOffer, err := p.Offers.Add(offerModel)
	if err != nil {
		return newOffer, err
	}
	query := fmt.Sprintf("INSERT INTO %s(job_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err = p.Db.Pool.Exec(context.Background(), query, newOffer.Model().Id, p.Relation.RelationId)
	if err != nil {
		return newOffer, err
	}
	return newOffer, nil
}
