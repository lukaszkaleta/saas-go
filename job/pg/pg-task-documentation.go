package pgjob

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgTaskDocumentation struct {
	db     *pg.PgDb
	taskId int64
}

func (pg *PgTaskDocumentation) Create(ctx context.Context, summary *universal.DescriptionModel) (job.TaskDocumentationEntry, error) {
	query := `
		insert into task_documentation_entry(task_id, summary_value, summary_image_url, action_created_by_id)
		values (@taskId, @summaryValue, @summaryImageUrl, @createdById)
		returning ` + TaskDocumentationEntryColumnString() + `
	`
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{
		"taskId":          pg.taskId,
		"summaryValue":    summary.Value,
		"summaryImageUrl": summary.ImageUrl,
		"createdById":     universal.CurrentUserId(ctx),
	})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapTaskDocumentationEntry(pg.db))
}

func (pg *PgTaskDocumentation) EntriesModels(ctx context.Context) ([]job.TaskDocumentationEntryModel, error) {
	query := TaskDocumentationEntrySelect() + " where task_id = @taskId order by action_created_at desc "
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"taskId": pg.taskId})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, MapTaskDocumentationEntryModel)
	if err != nil {
		return nil, err
	}
	entries := make([]job.TaskDocumentationEntryModel, len(models))
	for i := range models {
		entries[i] = *models[i]
		pgEntry := &PgTaskDocumentationEntry{db: pg.db, id: entries[i].Id, taskId: entries[i].TaskId}
		images, err := pgEntry.Images(ctx)
		if err == nil {
			entries[i].Images = images
		}
	}
	return entries, nil
}

func (pg *PgTaskDocumentation) assertAccess(ctx context.Context) (int64, error) {
	currentUser := universal.CurrentUserId(ctx)
	if currentUser == nil || *currentUser <= 0 {
		return 0, job.ErrTaskDocumentationMissingUser
	}
	var taskUserId int64
	var jobId int64
	var jobOwnerId int64
	query := `
		select t.user_id, t.job_id, j.action_created_by_id
		from task t
		join job j on j.id = t.job_id
		where t.id = @taskId
	`
	err := pg.db.Pool.QueryRow(ctx, query, pgx.NamedArgs{"taskId": pg.taskId}).Scan(&taskUserId, &jobId, &jobOwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, job.ErrTaskDocumentationTaskNotFound
		}
		return 0, err
	}
	if *currentUser != taskUserId && *currentUser != jobOwnerId {
		return 0, job.ErrTaskDocumentationAccessDenied
	}
	return jobId, nil
}
