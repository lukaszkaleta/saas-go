package filestore

import "github.com/lukaszkaleta/saas-go/universal"

// API

type FileSystem interface {
	Model() *FileSystemModel
	Update(newModel *FileSystemModel) error
	Records() Records
}

// Builder

// Model

type FileSystemModel struct {
	Id   int64 `json:"id"`
	Name *universal.NameModel
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

func (addr SolidFileSystem) Update(newModel *FileSystemModel) error {
	addr.model.Change(newModel)
	if addr.FileSystem == nil {
		return nil
	}
	return addr.FileSystem.Update(newModel)
}

func (addr SolidFileSystem) Model() *FileSystemModel {
	return addr.model
}

func (addr SolidFileSystem) Records() Records {
	if addr.FileSystem == nil {
		return NoRecords{}
	}
	return addr.FileSystem.Records()
}
