package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgTasks struct {
	db     *pg.PgDb
	UserId int64
}

func (pgTasks *PgTasks) Archived(ctx context.Context) ([]job.Task, error) {
	query := "select * from task where user_id = @userId and action_finished_at is not null and action_pay_at is not null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) WaitingForPayment(ctx context.Context) ([]job.Task, error) {
	query := "select * from task where user_id = @userId and action_finished_at is not null and action_pay_at is null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) InProgress(ctx context.Context) ([]job.Task, error) {
	query := "select * from task where user_id = @userId and action_finished_at is null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) Create(ctx context.Context, model *job.TaskModel) (job.Task, error) {
	currentUserId := universal.CurrentUserId(ctx)
	query := "INSERT INTO task(job_id, user_id, offer_id, action_created_by_id) values (@jobId, @userId, @offerId, @currentUserId) returning id"
	row := pgTasks.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"jobId":         model.JobId,
		"userId":        pgTasks.UserId,
		"offerId":       model.OfferId,
		"currentUserId": currentUserId,
	})
	err := row.Scan(&model.Id)
	if err != nil {
		return nil, err
	}
	return job.NewSolidTask(model, &PgTask{db: pgTasks.db, Id: model.Id}), nil
}

func NewPgTasks(db *pg.PgDb, userId int64) *PgTasks {
	return &PgTasks{db: db, UserId: userId}
}
