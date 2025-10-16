package pgoffer

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	pgFilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/offer"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgOffer struct {
	Db *pg.PgDb
	Id int64
}

func (pgOffer *PgOffer) Model() *offer.OfferModel {
	//TODO implement me
	panic("implement me")
}

func (pgOffer *PgOffer) Address() universal.Address {
	return &pgUniversal.PgAddress{pgOffer.Db, pgOffer.tableEntity()}
}

func (pgOffer *PgOffer) Position() universal.Position {
	return &pgUniversal.PgPosition{pgOffer.Db, pgOffer.tableEntity()}
}

func (pgOffer *PgOffer) Price() universal.Price {
	return &pgUniversal.PgPrice{pgOffer.Db, pgOffer.tableEntity()}
}

func (pgOffer *PgOffer) Description() universal.Description {
	return pgUniversal.NewPgDescriptionFromTable(pgOffer.Db, pgOffer.tableEntity())
}

func (pgOffer *PgOffer) FileSystem() filestore.FileSystem {
	return &pgFilestore.PgFileSystem{
		Db: pgOffer.Db,
		Owner: pg.RelationEntity{
			RelationId: pgOffer.Id,
			TableName:  "offer_filesystem",
			ColumnName: "offer_id",
		},
	}
}

func (pgOffer *PgOffer) State() universal.State {
	return pgUniversal.NewPgTimestampState(
		pgOffer.Db,
		pgOffer.tableEntity(),
		offer.OfferStatuses())
}

func (pgOffer *PgOffer) tableEntity() pg.TableEntity {
	return pgOffer.Db.TableEntity("job", pgOffer.Id)
}

func (pgOffer *PgOffer) localizationRelationEntity() pg.TableEntity {
	return pgOffer.Db.TableEntity("job", pgOffer.Id)
}

func MapOffer(row pgx.CollectableRow) (*offer.OfferModel, error) {
	offerModel := offer.EmptyOfferModel()

	nullTimePublished := sql.NullTime{}
	nullTimeClosed := sql.NullTime{}
	err := row.Scan(
		&offerModel.Id,
		&offerModel.Description.Value,
		&offerModel.Description.ImageUrl,
		&offerModel.Address.Line1,
		&offerModel.Address.Line2,
		&offerModel.Address.City,
		&offerModel.Address.PostalCode,
		&offerModel.Address.District,
		&offerModel.Position.Lat,
		&offerModel.Position.Lon,
		&offerModel.Price.Value,
		&offerModel.Price.Currency,
		&offerModel.State.Draft,
		&nullTimePublished,
		&nullTimeClosed,
	)
	offerModel.State.Published = nullTimePublished.Time
	offerModel.State.Closed = nullTimeClosed.Time
	if err != nil {
		return nil, err
	}
	return offerModel, nil
}
