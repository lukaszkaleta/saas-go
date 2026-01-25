package pgjob

import (
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	pgFilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/messages"
	pgMessages "github.com/lukaszkaleta/saas-go/messages/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

type PgJob struct {
	db *pg.PgDb
	Id int64
}

func (pgJob *PgJob) ID() int64 {
	return pgJob.Id
}

func (pgJob *PgJob) Model() *job.JobModel {
	//TODO implement me
	panic("implement me")
}

func (pgJob *PgJob) Address() universal.Address {
	return &pgUniversal.PgAddress{pgJob.db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Position() universal.Position {
	return &pgUniversal.PgPosition{pgJob.db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Price() universal.Price {
	return &pgUniversal.PgPrice{pgJob.db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Description() universal.Description {
	return pgUniversal.NewPgDescriptionFromTable(pgJob.db, pgJob.tableEntity())
}

func (pgJob *PgJob) FileSystem() filestore.FileSystem {
	return &pgFilestore.PgFileSystem{
		Db: pgJob.db,
		Owner: pg.RelationEntity{
			RelationId: pgJob.Id,
			TableName:  "job_filesystem",
			ColumnName: "job_id",
		},
	}
}

func (pgJob *PgJob) State() universal.State {
	return pgUniversal.NewPgTimestampState(
		pgJob.db,
		pgJob.tableEntity(),
		job.JobStatuses())
}

func (pgJob *PgJob) Actions() universal.Actions {
	return pgUniversal.NewPgActions(pgJob.db, pgJob.tableEntity())
}

func (pgJob *PgJob) Offers() job.Offers {
	return &PgOffers{db: pgJob.db, JobId: pgJob.Id}
}

func (pgJob *PgJob) Messages() messages.Messages {
	return pgMessages.NewPgMessages(pgJob.db, pg.TableEntity{Name: "job_message", Id: pgJob.Id})
}

func (pgJob *PgJob) tableEntity() pg.TableEntity {
	return pgJob.db.TableEntity("job", pgJob.Id)
}

func (pgJob *PgJob) localizationRelationEntity() pg.TableEntity {
	return pgJob.db.TableEntity("job", pgJob.Id)
}

// Mapping Job

func MapJob(db *pg.PgDb) pgx.RowToFunc[job.Job] {
	return func(row pgx.CollectableRow) (job.Job, error) {
		model, err := MapJobModel(row)
		if err != nil {
			return nil, err
		}
		pgJob := &PgJob{db: db, Id: model.Id}
		return job.NewSolidJob(model, pgJob), nil
	}
}

func MapJobModel(row pgx.CollectableRow) (*job.JobModel, error) {
	jobModel := job.EmptyJobModel()

	nullTimePublished := sql.NullTime{}
	nullTimeOccupied := sql.NullTime{}
	nullTimeClosed := sql.NullTime{}

	actionCreated := universal.ActionModel{Name: "created"}
	actions := make(map[string]*universal.ActionModel)
	actions["created"] = &actionCreated

	err := row.Scan(
		&jobModel.Id,
		&jobModel.Description.Value,
		&jobModel.Description.ImageUrl,
		&jobModel.Address.Line1,
		&jobModel.Address.Line2,
		&jobModel.Address.City,
		&jobModel.Address.PostalCode,
		&jobModel.Address.District,
		&jobModel.Position.Lat,
		&jobModel.Position.Lon,
		&jobModel.Price.Value,
		&jobModel.Price.Currency,
		&jobModel.Rating,
		&jobModel.State.Draft,
		&nullTimePublished,
		&nullTimeOccupied,
		&nullTimeClosed,
		&jobModel.Tags,
		&actionCreated.ById,
		&actionCreated.MadeAt,
	)
	jobModel.State.Published = nullTimePublished.Time
	jobModel.State.Occupied = nullTimeOccupied.Time
	jobModel.State.Closed = nullTimeClosed.Time
	jobModel.Actions = universal.ActionsModel{List: actions}
	if err != nil {
		return nil, err
	}
	return jobModel, nil
}

func JobColumns() []string {
	return []string{
		"id",
		"description_value",
		"description_image_url",
		"address_line_1",
		"address_line_2",
		"address_city",
		"address_postal_code",
		"address_district",
		"position_latitude",
		"position_longitude",
		"price_value",
		"price_currency",
		"rating",
		"status_draft",
		"status_published",
		"status_occupied",
		"status_closed",
		"tags",
		"action_created_by_id",
		"action_created_at",
	}
}

func JobColumnString() string {
	return strings.Join(JobColumns(), ",")
}

func JobSelect() string {
	return JobColumnsSelect() + " from job "
}

func JobColumnsSelect() string {
	return "select " + JobColumnString()
}

// Mapping search

func MapSearchJob() pgx.RowToFunc[*job.JobSearchOutput] {
	return func(row pgx.CollectableRow) (*job.JobSearchOutput, error) {
		jobModel := job.EmptyJobModel()

		nullTimePublished := sql.NullTime{}
		nullTimeOccupied := sql.NullTime{}
		nullTimeClosed := sql.NullTime{}

		actionCreated := universal.ActionModel{Name: "created"}
		actions := make(map[string]*universal.ActionModel)
		actions["created"] = &actionCreated

		jobSearchRanking := &job.JobSearchRanking{}
		jobSearchPaging := &job.JobSearchPaging{}
		jobSearchOutput := &job.JobSearchOutput{
			Model:   jobModel,
			Ranking: jobSearchRanking,
			Paging:  jobSearchPaging,
		}
		err := row.Scan(
			&jobSearchOutput.Ranking.Distance,
			&jobSearchOutput.Ranking.Rank,
			&jobModel.Id,
			&jobModel.Description.Value,
			&jobModel.Description.ImageUrl,
			&jobModel.Address.Line1,
			&jobModel.Address.Line2,
			&jobModel.Address.City,
			&jobModel.Address.PostalCode,
			&jobModel.Address.District,
			&jobModel.Position.Lat,
			&jobModel.Position.Lon,
			&jobModel.Price.Value,
			&jobModel.Price.Currency,
			&jobModel.Rating,
			&jobModel.State.Draft,
			&nullTimePublished,
			&nullTimeOccupied,
			&nullTimeClosed,
			&jobModel.Tags,
			&actionCreated.ById,
			&actionCreated.MadeAt,
		)
		jobModel.State.Published = nullTimePublished.Time
		jobModel.State.Occupied = nullTimeOccupied.Time
		jobModel.State.Closed = nullTimeClosed.Time
		jobModel.Actions = universal.ActionsModel{List: actions}
		if err != nil {
			return nil, err
		}
		return jobSearchOutput, nil
	}
}
