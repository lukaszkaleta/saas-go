package job

import (
	"time"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Job interface {
	Model() *JobModel
	Address() universal.Address
	Position() universal.Position
	Price() universal.Price
	Description() universal.Description
	FileSystem() filestore.FileSystem
	State() universal.State
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
	if !o.Published.IsZero() {
		return JobClosed
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

// Model

type JobModel struct {
	Id          int64                       `json:"id"`
	Position    *universal.PositionModel    `json:"position"`
	Price       *universal.PriceModel       `json:"price"`
	Address     *universal.AddressModel     `json:"address"`
	Description *universal.DescriptionModel `json:"description"`
	State       JobStatus                   `json:"state"`
}

func (m *JobModel) Hint() *JobHint {
	return &JobHint{
		Id:    m.Id,
		Price: m.Price.UserFriendly(),
	}
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

func NewSolidJob(model *JobModel, Job Job, id int64) Job {
	return &SolidJob{
		id,
		model,
		Job,
	}
}
func (solidJob *SolidJob) Model() *JobModel {
	return solidJob.model
}

func (solidJob *SolidJob) Position() universal.Position {
	if solidJob.Job != nil {
		return universal.NewSolidPosition(
			solidJob.Model().Position,
			solidJob.Job.Position(),
		)
	}
	return universal.NewSolidPosition(solidJob.Model().Position, nil)
}

func (solidJob *SolidJob) Price() universal.Price {
	if solidJob.Job != nil {
		return universal.NewSolidPrice(
			solidJob.Model().Price,
			solidJob.Job.Price(),
		)
	}
	return universal.NewSolidPrice(solidJob.Model().Price, nil)
}

func (solidJob *SolidJob) Address() universal.Address {
	if solidJob.Job != nil {
		return universal.NewSolidAddress(
			solidJob.Model().Address,
			solidJob.Job.Address(),
		)
	}
	return universal.NewSolidAddress(solidJob.Model().Address, nil)
}
func (solidJob *SolidJob) Description() universal.Description {
	if solidJob.Job != nil {
		return universal.NewSolidDescription(
			solidJob.Model().Description,
			solidJob.Job.Description(),
		)
	}
	return universal.NewSolidDescription(solidJob.Model().Description, nil)
}

func (solidJob *SolidJob) FileSystem() filestore.FileSystem {
	return solidJob.Job.FileSystem()
}

func (solidJob *SolidJob) State() universal.State {
	available := JobStatuses()
	current := solidJob.Model().State.Current()
	if solidJob.Job != nil {
		return universal.NewSolidState(
			current,
			available,
			solidJob.Job.State())
	}
	return universal.NewSolidState(current, available, nil)
}
