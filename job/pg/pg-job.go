package pgjob

import (
	"database/sql"

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
	Db *pg.PgDb
	Id int64
}

func (pgJob *PgJob) Model() *job.JobModel {
	//TODO implement me
	panic("implement me")
}

func (pgJob *PgJob) Address() universal.Address {
	return &pgUniversal.PgAddress{pgJob.Db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Position() universal.Position {
	return &pgUniversal.PgPosition{pgJob.Db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Price() universal.Price {
	return &pgUniversal.PgPrice{pgJob.Db, pgJob.tableEntity()}
}

func (pgJob *PgJob) Description() universal.Description {
	return pgUniversal.NewPgDescriptionFromTable(pgJob.Db, pgJob.tableEntity())
}

func (pgJob *PgJob) FileSystem() filestore.FileSystem {
	return &pgFilestore.PgFileSystem{
		Db: pgJob.Db,
		Owner: pg.RelationEntity{
			RelationId: pgJob.Id,
			TableName:  "job_filesystem",
			ColumnName: "job_id",
		},
	}
}

func (pgJob *PgJob) State() universal.State {
	return pgUniversal.NewPgTimestampState(
		pgJob.Db,
		pgJob.tableEntity(),
		job.JobStatuses())
}

func (pgJob *PgJob) Actions() universal.Actions {
	return pgUniversal.NewPgActions(pgJob.Db, pgJob.tableEntity())
}

func (pgJob *PgJob) Offers() job.Offers {
	return &PgOffers{Db: pgJob.Db, JobId: pgJob.Id}
}

func (pgJob *PgJob) Messages() messages.Messages {
	return nil
}

func (pgJob *PgJob) tableEntity() pg.TableEntity {
	return pgJob.Db.TableEntity("job", pgJob.Id)
}

func (pgJob *PgJob) localizationRelationEntity() pg.TableEntity {
	return pgJob.Db.TableEntity("job", pgJob.Id)
}

func MapJob(row pgx.CollectableRow) (*job.JobModel, error) {
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
