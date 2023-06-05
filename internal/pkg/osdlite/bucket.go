package osdlite

import "github.com/google/uuid"

type bucket struct {
	id       string `json:"id"`
	name     string `json:"name"`
	owner    string `json:"owner"`
	basepath string `json:"basepath"`
}

func NewBucket(name, owner, basepath string) *bucket {
	return &bucket{
		id:       uuid.NewString(),
		owner:    owner,
		basepath: basepath,
		name:     name,
	}
}

func (b bucket) TableName() string {
	return "buckets"
}
