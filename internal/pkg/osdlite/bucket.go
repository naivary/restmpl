package osdlite

import (
	"time"

	"github.com/google/uuid"
)

type bucket struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Basepath  string `json:"basepath"`
}

func NewBucket(name, owner, basepath string) *bucket {
	return &bucket{
		ID:        uuid.NewString(),
		Name:      name,
		Owner:     owner,
		Basepath:  basepath,
		CreatedAt: time.Now().Unix(),
	}
}

func (b bucket) TableName() string {
	return "buckets"
}

func (b bucket) Create() error {
	return nil
}
