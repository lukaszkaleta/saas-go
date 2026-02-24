package pgjob

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	pgFilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/universal"
	pgUniversal "github.com/lukaszkaleta/saas-go/universal/pg"
)

// PgTask implements job.Task backed by Postgres
type PgTask struct {
	db *pg.PgDb
	Id int64
}

func (p *PgTask) Job(ctx context.Context) (job.Job, error) {
	query := JobSelect() + " where id = (select job_id from task where id = @id)"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": p.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapJob(p.db))
}

func (p *PgTask) ID() int64 {
	return p.Id
}

func (p *PgTask) FileSystem() filestore.FileSystem {
	return pgFilestore.NewPgFileSystem(
		p.db,
		pg.RelationEntity{
			RelationId: p.Id,
			TableName:  "task_filesystem",
			ColumnName: "task_id",
		},
	)
}

func (p *PgTask) Model(ctx context.Context) (*job.TaskModel, error) {
	query := "select * from task where id = @id"
	rows, err := p.db.Pool.Query(ctx, query, pgx.NamedArgs{"id": p.Id})
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(rows, MapTaskModel)
}

func (p *PgTask) Summary() universal.Description {
	return pgUniversal.NewPgDescription(
		p.db,
		p.tableEntity(),
		"summary_value",
		"summary_image_url",
	)
}

func (p *PgTask) tableEntity() pg.TableEntity {
	return p.db.TableEntity("task", p.Id)
}

// Mapping helpers

func MapTask(db *pg.PgDb) pgx.RowToFunc[job.Task] {
	return func(row pgx.CollectableRow) (job.Task, error) {
		model, err := MapTaskModel(row)
		if err != nil {
			return nil, err
		}
		pgTask := &PgTask{db: db, Id: model.Id}
		return job.NewSolidTask(model, pgTask), nil
	}
}

func MapTaskModel(row pgx.CollectableRow) (*job.TaskModel, error) {
	m := &job.TaskModel{
		Summary: universal.EmptyDescriptionModel(),
		Actions: universal.EmptyActionsModel(),
	}
	created := universal.EmptyCreatedActionModel()
	finished := universal.EmptyActionModel("finished")
	paid := universal.EmptyActionModel("paid")
	if err := row.Scan(
		&m.Id,
		&m.JobId,
		&m.OfferId,
		&m.UserId,
		&m.Summary.Value,
		&m.Summary.ImageUrl,
		&created.ById,
		&created.MadeAt,
		&finished.ById,
		&finished.MadeAt,
		&paid.ById,
		&paid.MadeAt,
	); err != nil {
		return nil, err
	}
	m.Actions.List[created.Name] = created
	m.Actions.List[finished.Name] = finished
	m.Actions.List[paid.Name] = paid
	return m, nil
}
