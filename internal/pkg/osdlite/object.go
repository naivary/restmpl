package osdlite

import (
	"time"

	"github.com/google/uuid"
)

type object struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
	Version      int       `json:"version"`
	Name         string    `json:"name"`
	Owner        string    `json:"owner"`
	BucketID     string    `json:"bucketId"`
	Payload      []byte    `json:"payload"`
	Tags         []tag     `json:"tags"`
}

func NewObject(bucketID string) *object {
	return &object{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
		BucketID:     bucketID,
		Version:      1,
	}
}

func (o *object) Touch() {
	o.LastModified = time.Now()
}

func (o object) TableName() string {
	return "objects"
}
