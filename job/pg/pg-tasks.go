package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pguser "github.com/lukaszkaleta/saas-go/user/pg"
)

type PgTasks struct {
	db     *pg.PgDb
	UserId int64
}

func (pgTasks *PgTasks) ById(ctx context.Context, id int64) (job.Task, error) {
	query := "select * from task where id = @id"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapTask(pgTasks.db))
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

func (pgTasks *PgTasks) Current(ctx context.Context) (*job.TasksResult, error) {
	query := "select * from task where user_id = @userId and action_pay_at is null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	tasks, err := pgx.CollectRows(rows, MapTask(pgTasks.db))
	if err != nil {
		return nil, err
	}

	// Read Jobs
	jobs, err := pgTasks.readJobs(ctx, tasks)
	if err != nil {
		return nil, err
	}

	// Read persons
	personModels, err := pgTasks.readPersons(ctx, jobs)
	if err != nil {
		return nil, err
	}

	return job.NewTasksResult(tasks, jobs, personModels), nil
}

func (pgTasks *PgTasks) Completed(ctx context.Context) (*job.TasksResult, error) {

	// Read Tasks
	query := "select * from task where user_id = @userId and action_finished_at is not null and action_pay_at is not null"
	rows, err := pgTasks.db.Pool.Query(ctx, query, pgx.NamedArgs{"userId": pgTasks.UserId})
	if err != nil {
		return nil, err
	}
	tasks, err := pgx.CollectRows(rows, MapTask(pgTasks.db))

	// Read Jobs
	jobs, err := pgTasks.readJobs(ctx, tasks)
	if err != nil {
		return nil, err
	}

	// Read persons
	personModels, err := pgTasks.readPersons(ctx, jobs)
	if err != nil {
		return nil, err
	}

	return job.NewTasksResult(tasks, jobs, personModels), nil
}

func (pgTasks *PgTasks) readPersons(ctx context.Context, jobs []job.Job) ([]*universal.PersonModel, error) {
	ids := make([]*int64, 0, len(jobs))
	for _, j := range jobs {
		model, err := j.Model(ctx)
		if err != nil {
			return nil, err
		}
		ids = append(ids, model.Actions.CreatedById())
	}
	userSearch := pguser.NewPgUserSearch(pgTasks.db)
	return userSearch.PersonModelsByIds(ctx, ids)
}

func (pgTasks *PgTasks) readJobs(ctx context.Context, tasks []job.Task) ([]job.Job, error) {
	jobIds := make([]int64, 0, len(tasks))
	for _, task := range tasks {
		model, err := task.Model(ctx)
		if err != nil {
			return nil, err
		}
		jobIds = append(jobIds, model.JobId)
	}
	pgJobs := NewPgJobs(pgTasks.db, jobIds)
	return pgJobs.List(ctx)
}

func (pgTasks *PgTasks) Earnings(ctx context.Context) (map[string]*universal.PriceModel, error) {
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

	earnings := make(map[string]*universal.PriceModel)
	for rows.Next() {
		var status string
		var amount int
		var currency string
		err := rows.Scan(&status, &amount, &currency)
		if err != nil {
			return nil, err
		}
		earnings[status] = &universal.PriceModel{
			Value:    amount,
			Currency: currency,
		}
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
