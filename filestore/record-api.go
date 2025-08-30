package filestore

import "github.com/lukaszkaleta/saas-go/universal"

// API

type Record interface {
	Model() *RecordModel
	Update(newModel *RecordModel) error
}

// Builder

// Model

type RecordModel struct {
	Id          int64 `json:"id"`
	Url         string
	Name        *universal.NameModel
	Description *universal.DescriptionModel
}

func (model *RecordModel) Change(newModel *RecordModel) {
	model.Url = newModel.Url
	model.Name.Change(newModel.Name.Value)
	model.Description.Change(newModel.Description)
}

func EmptyRecordModel() *RecordModel {
	return &RecordModel{
		Id:          0,
		Url:         "",
		Name:        universal.EmptyNameModel(),
		Description: universal.EmptyDescriptionModel(),
	}
}

// Solid

type SolidRecord struct {
	model  *RecordModel
	Record Record
}

func NewSolidRecord(model *RecordModel, record Record) Record {
	return &SolidRecord{model, record}
}

func (addr SolidRecord) Update(newModel *RecordModel) error {
	addr.model.Change(newModel)
	if addr.Record == nil {
		return nil
	}
	return addr.Record.Update(newModel)
}

func (addr SolidRecord) Model() *RecordModel {
	return addr.model
}
