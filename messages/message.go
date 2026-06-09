package messages

import (
	"context"

	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

type OLDMessage interface {
	universal.Idable
	filestore.FileSystemAware
	Model(ctx context.Context) (*OLDMessageModel, error)
	Acknowledge(ctx context.Context) error
}

type OLDMessageModel struct {
	Id             int64                   `json:"id"`
	OwnerId        int64                   `json:"ownerId"`
	RecipientId    int64                   `json:"recipientId"`
	Value          string                  `json:"value"`
	ValueGenerated bool                    `json:"generated"`
	Actions        *universal.ActionsModel `json:"actions"`
}

func EmptyModel() *OLDMessageModel {
	return EmptyOwnerModel(0)
}

func (m OLDMessageModel) ID() int64 {
	return m.Id
}

func (m OLDMessageModel) GetActions() *universal.ActionsModel {
	return m.Actions
}

func EmptyOwnerModel(ownerId int64) *OLDMessageModel {
	return &OLDMessageModel{
		Id:          0,
		OwnerId:     ownerId,
		RecipientId: 0,
		Value:       "",
		Actions:     universal.EmptyActionsModel(),
	}
}

// Solid

type OLDSolidMessage struct {
	Id      int64
	model   *OLDMessageModel
	message OLDMessage
}

func (m *OLDSolidMessage) Acknowledge(ctx context.Context) error {
	if m.message != nil {
		return m.message.Acknowledge(ctx)
	}
	return nil
}

func NewSolidMessage(model *OLDMessageModel, message OLDMessage, id int64) OLDMessage {
	return &OLDSolidMessage{
		Id:      id,
		model:   model,
		message: message,
	}
}

func (m *OLDSolidMessage) FileSystem() filestore.FileSystem {
	return m.message.FileSystem()
}

func (m *OLDSolidMessage) Model(ctx context.Context) (*OLDMessageModel, error) {
	return m.model, nil
}

func (m *OLDSolidMessage) ID() int64 {
	return m.Id
}
