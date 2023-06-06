package osdlite

import (
	"time"

	"github.com/google/uuid"
)

type bucket struct {
	ID        string
	CreatedAt int64
	Name      string
	Basepath  string
	Owner     string
}

func NewBucket(name, basepath, owner string) *bucket {
	return &bucket{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		Name:      name,
		Basepath:  basepath,
		Owner:     owner,
	}
}

func (b bucket) TableName() string {
	return "buckets"
}
