package job

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/messages"
	"github.com/lukaszkaleta/saas-go/payment"
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Job interface {
	universal.Idable
	filestore.FileSystemAware
	universal.ActionsAware
	universal.Closable
	Model(ctx context.Context) (*JobModel, error)
	Address() universal.Address
	Position() universal.Position
	Price() universal.Price
	Description() universal.Description
	State() universal.State
	Offers() Offers
	Messages() messages.Messages
	MakeTask(ctx context.Context, offerId int64) error
	Payments() payment.Payments
	Ratings() universal.Ratings
	AssertJobOwnerAccess(ctx context.Context) error
	PersonModel(ctx context.Context) (*universal.PersonModel, error)
}

type JobStatus struct {
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
	return ""
}

const (
	JobPublished string = "published"
	JobOccupied  string = "occupied"
	JobClosed    string = "closed"
)

func Statuses() []string {
	return []string{JobPublished, JobOccupied, JobClosed}
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

func (m JobModel) Hint() *JobHint {
	return &JobHint{
		Id:    m.Id,
		Price: m.Price.UserFriendly(),
	}
}

func (m JobModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

func (m JobModel) ID() int64 {
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

func (solidJob *SolidJob) AssertJobOwnerAccess(ctx context.Context) error {
	currentUser := universal.CurrentUserId(ctx)
	if currentUser == nil || *currentUser <= 0 {
		return ErrTaskDocumentationMissingUser
	}
	model, err := solidJob.Model(ctx)
	if err != nil {
		return err
	}
	ownerId := model.Actions.CreatedById()
	if ownerId == nil || currentUser == nil || *ownerId != *currentUser {
		return ErrTaskDocumentationAccessDenied
	}
	return nil
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
	available := Statuses()
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

func (solidJob *SolidJob) Payments() payment.Payments {
	if solidJob.Job == nil {
		return nil
	}
	return solidJob.Job.Payments()
}

func (solidJob *SolidJob) Ratings() universal.Ratings {
	if solidJob.Job == nil {
		return universal.DummyRatings{}
	}
	return universal.NewSolidRatings(solidJob.Job.Ratings())
}

func (solidJob *SolidJob) MakeTask(ctx context.Context, offerId int64) error {
	return solidJob.Job.MakeTask(ctx, offerId)
}

func (solidJob *SolidJob) Close(ctx context.Context) error {
	return solidJob.Job.Close(ctx)
}

func (solidJob *SolidJob) Closed(ctx context.Context) (bool, error) {
	return solidJob.Job.Closed(ctx)
}

func (solidJob *SolidJob) PersonModel(ctx context.Context) (*universal.PersonModel, error) {
	return solidJob.Job.PersonModel(ctx)
}
