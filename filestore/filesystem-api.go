package filestore

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type FileSystem interface {
	universal.Idable
	Model(ctx context.Context) (*FileSystemModel, error)
	Update(ctx context.Context, newModel *FileSystemModel) error
	Records() Records
	Init(ctx context.Context) (int64, error)
	CheckExistence(ctx context.Context) (int64, error)
}

// Builder

// Model

type FileSystemModel struct {
	universal.Idable
	Id   int64 `json:"id"`
	Name *universal.NameModel
}

func (model *FileSystemModel) ID() int64 {
	return model.Id
}

func (model *FileSystemModel) Change(newModel *FileSystemModel) {
	model.Name.Change(newModel.Name.Value)
}

func EmptyFileSystemModel() *FileSystemModel {
	return &FileSystemModel{
		Id:   0,
		Name: universal.EmptyNameModel(),
	}
}

// Solid

type SolidFileSystem struct {
	model      *FileSystemModel
	FileSystem FileSystem
}

func NewSolidFileSystem(model *FileSystemModel, FileSystem FileSystem) FileSystem {
	return &SolidFileSystem{model, FileSystem}
}

func (addr SolidFileSystem) Update(ctx context.Context, newModel *FileSystemModel) error {
	addr.model.Change(newModel)
	if addr.FileSystem == nil {
		return nil
	}
	return addr.FileSystem.Update(ctx, newModel)
}

func (addr SolidFileSystem) Model(ctx context.Context) (*FileSystemModel, error) {
	return addr.model, nil
}

func (addr SolidFileSystem) Records() Records {
	if addr.FileSystem == nil {
		return NoRecords{}
	}
	return addr.FileSystem.Records()
}

func (addr SolidFileSystem) Init(ctx context.Context) (int64, error) {
	if addr.FileSystem == nil {
		return 0, nil
	}
	return addr.FileSystem.Init(ctx)
}

func (addr SolidFileSystem) CheckExistence(ctx context.Context) (int64, error) {
	if addr.FileSystem == nil {
		return 0, nil
	}
	return addr.FileSystem.CheckExistence(ctx)
}

func (addr SolidFileSystem) ID() int64 {
	return addr.model.Id
}
