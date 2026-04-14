package pgjob

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	pgFilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/messages"
	pgMessages "github.com/lukaszkaleta/saas-go/messages/pg"
	"github.com/lukaszkaleta/saas-go/payment"
	pgPayment "github.com/lukaszkaleta/saas-go/payment/pg"
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

func (pgJob *PgJob) OwnerUserId(ctx context.Context) (*int64, error) {
	model, err := pgJob.Actions().Model(ctx)
	if err != nil {
		return nil, err
	}
	return model.Created().ById, nil
}

func (pgJob *PgJob) Model(ctx context.Context) (*job.JobModel, error) {
	query := JobSelect() + " where id = @id"
	rows, err := pgJob.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapJobModel)
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
	return pgFilestore.NewPgFileSystem(
		pgJob.db,
		pg.RelationEntity{
			RelationId: pgJob.Id,
			TableName:  "job_filesystem",
			ColumnName: "job_id",
		},
	)
}

func (pgJob *PgJob) State() universal.State {
	return pgUniversal.NewPgTimestampState(
		pgJob.db,
		pgJob.tableEntity(),
		job.Statuses())
}

func (pgJob *PgJob) Actions() universal.Actions {
	return pgUniversal.NewPgActions(pgJob.db, pgJob.tableEntity())
}

func (pgJob *PgJob) Offers() job.Offers {
	return &PgOffers{db: pgJob.db, JobId: pgJob.Id}
}

func (pgJob *PgJob) Messages() messages.Messages {
	return pgMessages.NewPgMessages(
		pgJob.db,
		pg.RelationEntity{TableName: "job_message", ColumnName: "job_id", RelationId: pgJob.Id},
	)
}

func (pgJob *PgJob) Payments() payment.Payments {
	return pgPayment.NewPgPayments(pgJob.db, pgJob)
}

func (pgJob *PgJob) Ratings() universal.Ratings {
	return pgUniversal.NewPgRatings(pgJob.db, pgJob.tableEntity())
}

func (pgJob *PgJob) MakeTask(ctx context.Context, offerId int64) error {
	offer, err := pgJob.Offers().ById(ctx, offerId)
	if err != nil {
		return err
	}
	userId, err := universal.CreatedById[job.OfferModel](ctx, offer)
	if err != nil {
		return err
	}
	model := &job.TaskModel{UserId: userId, JobId: pgJob.Id, OfferId: offer.ID()}
	_, err = NewPgTasks(pgJob.db, userId).Create(ctx, model)
	if err != nil {
		return err
	}
	return nil
}

func (pgJob *PgJob) Close(ctx context.Context) error {
	err := pgJob.State().Change(ctx, job.JobClosed)
	if err != nil {
		return err
	}
	acceptedOffer, err := pgJob.Offers().Accepted(ctx)
	if err != nil {
		return err
	}
	userId, err := universal.CreatedById[job.OfferModel](ctx, acceptedOffer)
	if err != nil {
		return err
	}
	task, err := NewPgTasks(pgJob.db, userId).ByJobId(ctx, pgJob.Id)
	if err != nil {
		return err
	}
	err = task.Finish(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (pgJob *PgJob) Closed(ctx context.Context) (bool, error) {
	name, err := pgJob.State().Name(ctx)
	if err != nil {
		return false, err
	}
	return name == "closed", nil
}

func (pgJob *PgJob) PersonModel(ctx context.Context) (*universal.PersonModel, error) {
	query := pgUniversal.PersonColumnsSelectWithPrefix("u") + `
		from job j
		join users u on j.action_created_by_id = u.id
		where j.id = $1`

	rows, _ := pgJob.db.Pool.Query(ctx, query, pgJob.Id)
	return pgx.CollectOneRow(rows, pgUniversal.MapPersonModel)
}

func (pgJob *PgJob) tableEntity() pg.TableEntity {
	return pgJob.db.TableEntity("job", pgJob.Id)
}

func (pgJob *PgJob) localizationRelationEntity() pg.TableEntity {
	return pgJob.db.TableEntity("job", pgJob.Id)
}

func (pgJob *PgJob) AssertJobOwnerAccess(ctx context.Context) error {
	currentUser := universal.CurrentUserId(ctx)
	if currentUser == nil || *currentUser <= 0 {
		return job.ErrTaskDocumentationMissingUser
	}
	ownerId := int64(0)
	query := "select action_created_by_id from job where id = $1"
	if err := pgJob.db.Pool.QueryRow(ctx, query, pgJob.Id).Scan(&ownerId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return job.ErrTaskDocumentationJobNotFound
		}
		return err
	}
	if ownerId != *currentUser {
		return job.ErrTaskDocumentationAccessDenied
	}
	return nil
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

	actionCreatedModel := universal.EmptyCreatedActionModel()
	actions := make(map[string]*universal.ActionModel)
	actions[actionCreatedModel.Name] = actionCreatedModel

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
		&nullTimePublished,
		&nullTimeOccupied,
		&nullTimeClosed,
		&jobModel.Tags,
		&actionCreatedModel.ById,
		&actionCreatedModel.MadeAt,
	)
	jobModel.State.Published = nullTimePublished.Time
	jobModel.State.Occupied = nullTimeOccupied.Time
	jobModel.State.Closed = nullTimeClosed.Time
	jobModel.Actions = &universal.ActionsModel{List: actions}
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
		"status_published",
		"status_occupied",
		"status_closed",
		"tags",
		"action_created_by_id",
		"action_created_at",
	}
}

func MapJobColumns(mapper func(column string) string) []string {
	originalColumns := JobColumns()
	columns := make([]string, len(originalColumns))
	for i := range originalColumns {
		columns[i] = mapper(originalColumns[i])
	}
	return columns
}

func JobColumnString() string {
	return strings.Join(JobColumns(), ",")
}

func MapJobColumnString(mapper func(column string) string) string {
	return strings.Join(MapJobColumns(mapper), ",")
}

func JobSelect() string {
	return JobColumnsSelect() + " from job "
}

func JobColumnsSelect() string {
	return "select " + JobColumnString()
}

func MapJobColumnsSelect(mapper func(column string) string) string {
	return "select " + MapJobColumnString(mapper)
}

func JobColumnsSelectWithPrefix(prefix string) string {
	return MapJobColumnsSelect(
		func(c string) string {
			return prefix + "." + c
		},
	)
}

// Mapping search

func MapSearchJob() pgx.RowToFunc[*job.JobSearchResult] {
	return func(row pgx.CollectableRow) (*job.JobSearchResult, error) {
		jobModel := job.EmptyJobModel()

		nullTimePublished := sql.NullTime{}
		nullTimeOccupied := sql.NullTime{}
		nullTimeClosed := sql.NullTime{}

		actionCreated := universal.ActionModel{Name: "created"}
		actions := make(map[string]*universal.ActionModel)
		actions["created"] = &actionCreated

		jobSearchRanking := &job.JobSearchRanking{}
		jobSearchPaging := &job.JobSearchPaging{}
		jobSearchOutput := &job.JobSearchResult{
			Model:   jobModel,
			Person:  universal.EmptyPersonModel(),
			Ranking: jobSearchRanking,
			Paging:  jobSearchPaging,
		}
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
			&nullTimePublished,
			&nullTimeOccupied,
			&nullTimeClosed,
			&jobModel.Tags,
			&actionCreated.ById,
			&actionCreated.MadeAt,

			&jobSearchOutput.Ranking.Distance,
			&jobSearchOutput.Ranking.Rank,
		)
		jobModel.State.Published = nullTimePublished.Time
		jobModel.State.Occupied = nullTimeOccupied.Time
		jobModel.State.Closed = nullTimeClosed.Time
		jobModel.Actions = &universal.ActionsModel{List: actions}
		if err != nil {
			return nil, err
		}
		return jobSearchOutput, nil
	}
}
