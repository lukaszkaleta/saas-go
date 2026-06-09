package pg

import "github.com/lukaszkaleta/saas-go/database/pg"

type OLDPgRelationMessages struct {
	db       *pg.PgDb
	Messages *OLDPgMessages
	relation pg.RelationEntity
}

func NewOLDPgRelationMessages(pgMessages *OLDPgMessages, relation pg.RelationEntity) OLDPgRelationMessages {
	return OLDPgRelationMessages{
		db:       pgMessages.db,
		relation: relation,
		Messages: pgMessages,
	}
}
