package pgjob

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgTaskDocumentationEntry struct {
	db     *pg.PgDb
	id     int64
	taskId int64
}

func NewPgTaskDocumentationEntry(db *pg.PgDb, taskId int64) job.TaskDocumentation {
	return PgTaskDocumentation{taskId: taskId, db: db}
}

func (pg *PgTaskDocumentationEntry) ID() int64 {
	return pg.id
}

func (pg *PgTaskDocumentationEntry) Model(ctx context.Context) (*job.TaskDocumentationEntryModel, error) {
	query := TaskDocumentationEntrySelect() + " where id = @id "
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": pg.id})
	if err != nil {
		return nil, err
	}
	model, err := pgx.CollectOneRow(rows, MapTaskDocumentationEntryModel)
	if err != nil {
		return nil, err
	}
	images, err := pg.Images(ctx)
	if err != nil {
		model.Images = []string{}
	} else {
		model.Images = images
	}
	return model, nil
}

func (pg *PgTaskDocumentationEntry) Images(ctx context.Context) ([]string, error) {
	query := `
		select fr.description_image_url
		from filestore_record fr
		where fr.id in (
			select fsr.record_id
			from filesystem_record fsr
			where fsr.filesystem_id in (
				select tef.filesystem_id
				from task_documentation_entry_filesystem tef
				where tef.entry_id = @entryId
			)
		)
		and fr.description_image_url <> ''
	`
	rows, err := pg.db.Pool.Query(ctx, query, pgx.NamedArgs{"entryId": pg.id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func MapTaskDocumentationEntry(db *pg.PgDb) pgx.RowToFunc[job.TaskDocumentationEntry] {
	return func(row pgx.CollectableRow) (job.TaskDocumentationEntry, error) {
		model, err := MapTaskDocumentationEntryModel(row)
		if err != nil {
			return nil, err
		}
		pgEntry := &PgTaskDocumentationEntry{db: db, id: model.Id, taskId: model.TaskId}
		return job.NewSolidTaskDocumentationEntry(model, pgEntry), nil
	}
}

func MapTaskDocumentationEntryModel(row pgx.CollectableRow) (*job.TaskDocumentationEntryModel, error) {
	var entry job.TaskDocumentationEntryModel
	summary := universal.EmptyDescriptionModel()
	if err := row.Scan(
		&entry.Id,
		&entry.TaskId,
		&summary.Value,
		&summary.ImageUrl,
		&entry.CreatedById,
		&entry.CreatedAt,
	); err != nil {
		return nil, err
	}
	entry.Summary = summary
	return &entry, nil
}

func TaskDocumentationEntryColumns() []string {
	return []string{
		"id",
		"task_id",
		"summary_value",
		"summary_image_url",
		"action_created_by_id",
		"action_created_at",
	}
}

func MapTaskDocumentationEntryColumns(mapper func(column string) string) []string {
	originalColumns := TaskDocumentationEntryColumns()
	columns := make([]string, len(originalColumns))
	for i := range originalColumns {
		columns[i] = mapper(originalColumns[i])
	}
	return columns
}

func TaskDocumentationEntryColumnString() string {
	return strings.Join(TaskDocumentationEntryColumns(), ",")
}

func MapTaskDocumentationEntryColumnString(mapper func(column string) string) string {
	return strings.Join(MapTaskDocumentationEntryColumns(mapper), ",")
}

func TaskDocumentationEntrySelect() string {
	return TaskDocumentationEntryColumnsSelect() + " from task_documentation_entry "
}

func TaskDocumentationEntryColumnsSelect() string {
	return "select " + TaskDocumentationEntryColumnString()
}

func MapTaskDocumentationEntryColumnsSelect(mapper func(column string) string) string {
	return "select " + MapTaskDocumentationEntryColumnString(mapper)
}

func TaskDocumentationEntryColumnsSelectWithPrefix(prefix string) string {
	return MapTaskDocumentationEntryColumnsSelect(
		func(c string) string {
			return prefix + "." + c
		},
	)
}
