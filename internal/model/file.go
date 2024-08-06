package model

import "time"

type File struct {
	Name         string
	Key          string
	Size         int64
	LastModified time.Time
}
