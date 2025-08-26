package filestore

type FileSystems interface {
	Add(name string, ownerId int64) (FileSystem, error)
}
