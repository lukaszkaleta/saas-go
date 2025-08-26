package universal

import (
	"context"
	"fmt"
	"github.com/lukaszkaleta/saas-go/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgLocalizations struct {
	Db    *pg.PgDb
	Owner *pg.RelationEntity
}

func (pgLocalizations *PgLocalizations) Add(country string, language string, translation string) (universal.Localization, error) {
	slugedName := universal.SluggedName(translation)
	localizationId := int64(0)

	query := fmt.Sprintf("INSERT INTO %s (%s, country, language, translation_value, translation_slug) VALUES( $1, $2, $3, $4, $5) returning id",
		pgLocalizations.Owner.TableName,
		pgLocalizations.Owner.ColumnName,
	)
	rows, err := pgLocalizations.Db.Pool.Query(
		context.Background(),
		query,
		pgLocalizations.Owner.RelationId,
		country,
		language,
		slugedName.Value,
		slugedName.Slug)
	if err != nil {
		return nil, err
	}
	rows.Scan(&localizationId)
	rows.Close()
	pgLocalization := PgLocalization{
		Db:    pgLocalizations.Db,
		Id:    localizationId,
		Owner: pgLocalizations.Owner,
	}
	return universal.NewSolidLocalization(
		&universal.LocalizationModel{
			Id:          localizationId,
			OwnerId:     pgLocalizations.Owner.RelationId,
			Country:     country,
			Language:    language,
			Translation: slugedName,
		},
		&pgLocalization,
	), nil
}
