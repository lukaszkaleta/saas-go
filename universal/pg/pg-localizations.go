package postgres

import (
	"context"
	"fmt"
	"naborly/internal/api/common"postgres2 "naborly/internal/postgres"

)

type PgLocalizations struct {
	Db    *postgres2.PgDb
	Owner *RelationEntity
}

func (pgLocalizations *PgLocalizations) Add(country string, language string, translation string) (common.Localization, error) {
	slugedName := common.SluggedName(translation)
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
	return common.NewSolidLocalization(
		&common.LocalizationModel{
			Id:          localizationId,
			OwnerId:     pgLocalizations.Owner.RelationId,
			Country:     country,
			Language:    language,
			Translation: slugedName,
		},
		&pgLocalization,
	), nil
}
