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

func (pgOffers *PgOffers) AddFromPosition(model *universal.PositionModel) (offer.Offer, error) {
	offerId := int64(0)
	query := "INSERT INTO offer(position_latitude, position_longitude) VALUES( $1, $2 ) returning id"
	row := pgOffers.Db.Pool.QueryRow(context.Background(), query, model.Lat, model.Lon)
	row.Scan(&offerId)
	pgOffer := PgOffer{
		Db: pgOffers.Db,
		Id: offerId,
	}
	return offer.NewSolidOffer(
		&offer.OfferModel{
			Id:          offerId,
			Description: &universal.DescriptionModel{},
			Position:    model,
			Address:     &universal.AddressModel{},
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

func (p PgRelationOffers) AddFromPosition(model *universal.PositionModel) (offer.Offer, error) {
	newOffer, err := p.Offers.AddFromPosition(model)
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
