package pg

import "github.com/lukaszkaleta/saas-go/database/pg"

type PgRelationMessages struct {
	Db       *pg.PgDb
	Messages *PgMessages
	relation pg.RelationEntity
}

func NewPgRelationMessages(pgMessages *PgMessages, relation pg.RelationEntity) PgRelationMessages {
	return PgRelationMessages{
		Db:       pgMessages.Db,
		relation: relation,
		Messages: pgMessages,
	}
}
