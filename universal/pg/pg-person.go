package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgPerson struct {
	Db          *pg.PgDb
	TableEntity pg.TableEntity
}

func (pgPerson *PgPerson) Ratings() universal.Ratings {
	return NewPgRatings(pgPerson.Db, pgPerson.TableEntity)
}

func (pgPerson *PgPerson) Update(ctx context.Context, model *universal.PersonModel) error {
	query := fmt.Sprintf("update %s set person_first_name = $1, person_last_name = $2, person_email = $3, person_phone = $4 where id = $5", pgPerson.TableEntity.Name)
	_, err := pgPerson.Db.Pool.Exec(ctx, query, model.FirstName, model.LastName, model.Email, model.Phone, pgPerson.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pgPerson *PgPerson) UpdateAverageRating(ctx context.Context, score int) error {
	query := fmt.Sprintf("update %s set ratings_average = $1 where id = $2", pgPerson.TableEntity.Name)
	_, err := pgPerson.Db.Pool.Exec(ctx, query, score, pgPerson.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (pgPerson *PgPerson) Model(_ context.Context) *universal.PersonModel {
	return &universal.PersonModel{}
}

func (pgPerson *PgPerson) ID() int64 {
	return pgPerson.TableEntity.Id
}

func (pgPerson *PgPerson) Avatar(_ context.Context) universal.Description {
	return NewPgDescription(
		pgPerson.Db,
		pgPerson.TableEntity,
		"avatar_description_value",
		"avatar_description_image_url")
}

func MapPerson(db *pg.PgDb, tableEntity pg.TableEntity) pgx.RowToFunc[universal.Person] {
	return func(row pgx.CollectableRow) (universal.Person, error) {
		model, err := MapPersonModel(row)
		if err != nil {
			return nil, err
		}
		pgPerson := &PgPerson{Db: db, TableEntity: tableEntity}
		return universal.NewSolidPerson(model, pgPerson), nil
	}
}

func MapPersonModel(row pgx.CollectableRow) (*universal.PersonModel, error) {
	personModel := universal.EmptyPersonModel()
	err := row.Scan(
		&personModel.Id,
		&personModel.FirstName,
		&personModel.LastName,
		&personModel.Email,
		&personModel.Phone,
		&personModel.Avatar.Value,
		&personModel.Avatar.ImageUrl,
		&personModel.AverageRating,
	)
	if err != nil {
		return nil, err
	}
	return personModel, nil
}

func PersonColumns() []string {
	return []string{
		"id",
		"person_first_name",
		"person_last_name",
		"person_email",
		"person_phone",
		"avatar_description_value",
		"avatar_description_image_url",
		"ratings_average",
	}
}

func PersonColumnString() string {
	return strings.Join(PersonColumns(), ",")
}

func PersonColumnsSelect() string {
	return "select " + PersonColumnString()
}

func MapPersonColumns(mapper func(column string) string) []string {
	originalColumns := PersonColumns()
	columns := make([]string, len(originalColumns))
	for i := range originalColumns {
		columns[i] = mapper(originalColumns[i])
	}
	return columns
}

func MapPersonColumnString(mapper func(column string) string) string {
	return strings.Join(MapPersonColumns(mapper), ",")
}

func MapPersonColumnsSelect(mapper func(column string) string) string {
	return "select " + MapPersonColumnString(mapper)
}

func PersonColumnsSelectWithPrefix(prefix string) string {
	return MapPersonColumnsSelect(
		func(c string) string {
			return prefix + "." + c
		},
	)
}
