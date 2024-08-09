package model

import (
	"time"

	"github.com/XxThunderBlastxX/neoshare/internal/repository"
)

type File struct {
	Name         string
	Key          string
	Size         int64
	LastModified time.Time
}

func (f *File) FromRepositoryFile(file repository.File) {
	f.Name = file.Name
	f.Key = file.Key
	f.Size = int64(file.Size)
	f.LastModified = file.LastModified
}
