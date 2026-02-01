package filestore

import (
	"context"
	"os"

	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Record interface {
	Description() universal.Description
	Model(ctx context.Context) *RecordModel
	Update(ctx context.Context, newModel *RecordModel) error
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

func FileRecordModel(file os.File) *RecordModel {
	return &RecordModel{
		Id:          0,
		Url:         file.Name(),
		Name:        universal.SluggedName(file.Name()),
		Description: universal.EmptyDescriptionModel(),
	}
}

// Solid

type SolidRecord struct {
	model  *RecordModel
	Record Record
}

func (record SolidRecord) Description() universal.Description {
	return record.Record.Description()
}

func NewSolidRecord(model *RecordModel, record Record) Record {
	return &SolidRecord{model, record}
}

func (record SolidRecord) Update(ctx context.Context, newModel *RecordModel) error {
	record.model.Change(newModel)
	if record.Record == nil {
		return nil
	}
	return record.Record.Update(ctx, newModel)
}

func (record SolidRecord) Model(ctx context.Context) *RecordModel {
	return record.model
}
