package pg

import "github.com/lukaszkaleta/saas-go/database/pg"

type PgRelationMessages struct {
	db       *pg.PgDb
	Messages *PgMessages
	relation pg.RelationEntity
}

func NewPgRelationMessages(pgMessages *PgMessages, relation pg.RelationEntity) PgRelationMessages {
	return PgRelationMessages{
		db:       pgMessages.db,
		relation: relation,
		Messages: pgMessages,
	}
}
