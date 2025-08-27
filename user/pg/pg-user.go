package postgres

import (
	"naborly/internal/api/common"
	"naborly/internal/api/offer"
	"naborly/internal/api/rating"
	"naborly/internal/api/user"postgres2 "naborly/internal/postgres"

)

type PgUser struct {
	Db *PgDb
	Id int64
}

func (pgUser *PgUser) Address() common.Address {
	return &PgAddress{pgUser.Db, pgUser.tableEntity()}
}

func (pgUser *PgUser) Model() *user.UserModel {
	return &user.UserModel{}
}

func (pgUser *PgUser) Person() common.Person {
	return &PgPerson{pgUser.Db, pgUser.tableEntity()}
}

func (pgUser *PgUser) Ratings() rating.Ratings {
	return NewPgRatings(pgUser.Db, pgUser.tableEntity())
}

func (pgUser *PgUser) Settings() user.UserSettings {
	return NewPgUserSettings(pgUser.Db, pgUser.Id)
}

func (pgUser *PgUser) Account() user.Account {
	return NewPgAccount(pgUser.Db, pgUser.Id)
}

func (pgUser *PgUser) Archive() error {
	return nil
}

func (pgUser *PgUser) Offers() offer.Offers {
	tableEntity := pgUser.tableEntity()
	return postgres2.PgRelationOffers{
		Db:       pgUser.Db,
		offers:   &postgres2.PgOffers{Db: pgUser.Db},
		relation: tableEntity.RelationEntityWithColumnName("user_offer", "user_id"),
	}
}

func (pgUser *PgUser) tableEntity() TableEntity {
	return pgUser.Db.tableEntity("users", pgUser.Id)
}
