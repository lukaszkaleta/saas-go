package universal

import (
	"context"
	"fmt"
	"naborly/internal/api/common"
	postgres2 "naborly/internal/postgres"
)

type PgLocalization struct {
	Db    *pg.PgDb
	Id    int64
	Owner *RelationEntity
}

func (p *PgLocalization) Update(model *universal.LocalizationModel) error {
	translation := model.Translation
	query := fmt.Sprintf("update %s set translation_value = $1, translation_slug = $2 where %s = $3 and language = $4 and country = $5", p.Owner.TableName, p.Owner.ColumnName)
	_, err := p.Db.Pool.Exec(context.Background(), query, translation.Value, translation.Slug, p.Owner.RelationId, model.Language, model.Country)
	if err != nil {
		return err
	}
	return nil
}

func (p *PgLocalization) Model() *universal.LocalizationModel {
	return &universal.LocalizationModel{}
}
