// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"time"
)

type File struct {
	Key          string
	Name         string
	Size         int32
	LastModified time.Time
}