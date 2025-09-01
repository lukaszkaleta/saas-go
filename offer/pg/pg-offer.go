package pgoffer

import (
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

func (pgOffer *PgOffer) tableEntity() pg.TableEntity {
	return pgOffer.Db.TableEntity("offer", pgOffer.Id)
}

func (pgOffer *PgOffer) localizationRelationEntity() pg.TableEntity {
	return pgOffer.Db.TableEntity("offer", pgOffer.Id)
}
