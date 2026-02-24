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

func CurrentUserTasks(db *pg.PgDb, ctx context.Context) job.Tasks {
	return &PgTasks{db: db, UserId: *universal.CurrentUserId(ctx)}
}

func (pgTasks *PgTasks) ByJobId(ctx context.Context, jobId int64) (job.Task, error) {
	query := "select * from task where user_id = @userId and job_id = @jobId"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId, "jobId": jobId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) Current(ctx context.Context) ([]job.Task, error) {
	query := "select * from task where user_id = @userId and action_pay_at is null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) Completed(ctx context.Context) ([]job.Task, error) {
	query := "select * from task where user_id = @userId and action_finished_at is not null and action_pay_at is not null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, MapTask(pgTasks.db))
}

func (pgTasks *PgTasks) Earnings(ctx context.Context) (map[string]universal.Price, error) {
	query := `
		SELECT
			CASE
				WHEN t.action_pay_at IS NOT NULL THEN 'completed'
				WHEN t.action_finished_at IS NOT NULL THEN 'awaitingPayment'
				WHEN t.action_created_at IS NOT NULL THEN 'inProgress'
				ELSE 'unknown'
				END AS task_status,
			SUM(o.price_value) AS total_amount,
			o.price_currency
		FROM task t
				 JOIN job_offer o ON o.id = t.offer_id
		WHERE t.user_id = @userId
		GROUP BY
			task_status,
			o.price_currency;
`
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	earnings := make(map[string]universal.Price)
	for rows.Next() {
		var status string
		var amount int
		var currency string
		err := rows.Scan(&status, &amount, &currency)
		if err != nil {
			return nil, err
		}
		earnings[status] = universal.PriceFromModel(&universal.PriceModel{
			Value:    amount,
			Currency: currency,
		})
	}
	return earnings, nil
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
