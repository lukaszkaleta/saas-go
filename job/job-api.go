package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Job interface {
	universal.Idable
	Model(ctx context.Context) (*JobModel, error)
	Address() universal.Address
	Position() universal.Position
	Price() universal.Price
	Description() universal.Description
	FileSystem() filestore.FileSystem
	State() universal.State
	Actions() universal.Actions
	Offers() Offers
	Messages() messages.Messages
}

type JobStatus struct {
	Draft     time.Time `json:"draft"`
	Published time.Time `json:"published"`
	Occupied  time.Time `json:"occupied"`
	Closed    time.Time `json:"closed"`
}

func (o *JobStatus) Current() string {
	if !o.Closed.IsZero() {
		return JobClosed
	}
	if !o.Occupied.IsZero() {
		return JobOccupied
	}
	if !o.Published.IsZero() {
		return JobPublished
	}
	return JobDraft
}

const (
	JobDraft     string = "draft"
	JobPublished string = "published"
	JobOccupied  string = "occupied"
	JobClosed    string = "closed"
)

func JobStatuses() []string {
	return []string{JobDraft, JobPublished, JobOccupied, JobClosed}
}

func PublicStatuses() []string {
	return []string{JobPublished}
}

// JobModel

type JobModel struct {
	Id          int64                       `json:"id"`
	Position    *universal.PositionModel    `json:"position"`
	Price       *universal.PriceModel       `json:"price"`
	Rating      int                         `json:"rating"`
	Address     *universal.AddressModel     `json:"address"`
	Description *universal.DescriptionModel `json:"description"`
	State       JobStatus                   `json:"state"`
	Tags        []string                    `json:"tags"`
	Actions     *universal.ActionsModel     `json:"actions"`
}

func (m *JobModel) Hint() *JobHint {
	return &JobHint{
		Id:    m.Id,
		Price: m.Price.UserFriendly(),
	}
}

func (m *JobModel) ID() int64 {
	return m.Id
}

func EmptyJobModel() *JobModel {
	return &JobModel{
		Id:          0,
		Position:    &universal.PositionModel{},
		Address:     &universal.AddressModel{},
		Description: universal.EmptyDescriptionModel(),
		Price:       &universal.PriceModel{},
	}
}

type JobHint struct {
	Id       int64                    `json:"id"`
	Position *universal.PositionModel `json:"position"`
	Price    string                   `json:"price"`
}

// Solid

type SolidJob struct {
	Id    int64
	model *JobModel
	Job   Job
}

func NewSolidJob(model *JobModel, Job Job) Job {
	return &SolidJob{
		model.Id,
		model,
		Job,
	}
}

func (solidJob *SolidJob) ID() int64 {
	return solidJob.Id
}

func (solidJob *SolidJob) Model(ctx context.Context) (*JobModel, error) {
	return solidJob.model, nil
}

func (solidJob *SolidJob) Position() universal.Position {
	if solidJob.Job != nil {
		return universal.NewSolidPosition(
			solidJob.model.Position,
			solidJob.Job.Position(),
		)
	}
	return universal.NewSolidPosition(solidJob.model.Position, nil)
}

func (solidJob *SolidJob) Price() universal.Price {
	if solidJob.Job != nil {
		return universal.NewSolidPrice(
			solidJob.model.Price,
			solidJob.Job.Price(),
		)
	}
	return universal.NewSolidPrice(solidJob.model.Price, nil)
}

func (solidJob *SolidJob) Address() universal.Address {
	if solidJob.Job != nil {
		return universal.NewSolidAddress(
			solidJob.model.Address,
			solidJob.Job.Address(),
		)
	}
	return universal.NewSolidAddress(solidJob.model.Address, nil)
}
func (solidJob *SolidJob) Description() universal.Description {
	if solidJob.Job != nil {
		return universal.NewSolidDescription(
			solidJob.model.Description,
			solidJob.Job.Description(),
		)
	}
	return universal.NewSolidDescription(solidJob.model.Description, nil)
}

func (solidJob *SolidJob) FileSystem() filestore.FileSystem {
	return solidJob.Job.FileSystem()
}

func (solidJob *SolidJob) State() universal.State {
	available := JobStatuses()
	current := solidJob.model.State.Current()
	if solidJob.Job != nil {
		return universal.NewSolidState(
			current,
			available,
			solidJob.Job.State())
	}
	return universal.NewSolidState(current, available, nil)
}

func (solidJob *SolidJob) Actions() universal.Actions {
	return solidJob.Job.Actions()
}

func (solidJob *SolidJob) Offers() Offers {
	if solidJob.Job == nil {
		return NoOffers{}
	}
	return solidJob.Job.Offers()
}

func (solidJob *SolidJob) Messages() messages.Messages {
	if solidJob.Job == nil {
		return nil
	}
	return solidJob.Job.Messages()
}
