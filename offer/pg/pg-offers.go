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

func (pgOffers *PgOffers) AddWithPlace(positionModel *universal.PositionModel, addressModel *universal.AddressModel) (offer.Offer, error) {
	offerId := int64(0)
	query := "INSERT INTO job (position_latitude, position_longitude, address_line_1, address_line_2, address_city, address_postal_code, address_district) VALUES( $1, $2, $3, $4, $5, $6, $7 ) returning id"
	row := pgOffers.Db.Pool.QueryRow(
		context.Background(),
		query,
		positionModel.Lat,
		positionModel.Lon,
		addressModel.Line1,
		addressModel.Line2,
		addressModel.City,
		addressModel.PostalCode,
		addressModel.District)
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
			Position:    positionModel,
			Address:     addressModel,
			Price:       &universal.PriceModel{},
			State:       offer.OfferStatus{Draft: time.Now()},
		},
		&pgOffer,
		offerId,
	), nil
}

// Relation

type PgRelationOffers struct {
	Db       *pg.PgDb
	Offers   *PgOffers
	Relation pg.RelationEntity
}

func NewPgRelationOffers(pfOffers *PgOffers, relation pg.RelationEntity) PgRelationOffers {
	return PgRelationOffers{
		Db:       pfOffers.Db,
		Offers:   pfOffers,
		Relation: relation,
	}
}

func (p PgRelationOffers) AddWithPlace(positionModel *universal.PositionModel, addressModel *universal.AddressModel) (offer.Offer, error) {
	newOffer, err := p.Offers.AddWithPlace(positionModel, addressModel)
	if err != nil {
		return newOffer, err
	}
	query := fmt.Sprintf("INSERT INTO %s(offer_id, %s) VALUES( $1, $2 )", p.Relation.TableName, p.Relation.ColumnName)
	_, err = p.Db.Pool.Exec(context.Background(), query, newOffer.Model().Id, p.Relation.RelationId)
	if err != nil {
		return newOffer, err
	}
	return newOffer, nil
}
