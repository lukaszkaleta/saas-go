package filestore

import "context"

type FileSystems interface {
	Add(ctx context.Context, name string, ownerId int64) (FileSystem, error)
}
