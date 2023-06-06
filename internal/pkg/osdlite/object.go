package osdlite

import (
	"bytes"
	"time"

	"github.com/google/uuid"
)

type object struct {
	ID           string       `json:"id"`
	CreatedAt    int64        `json:"createdAt"`
	LastModified int64        `json:"lastModified"`
	Version      int          `json:"version"`
	Name         string       `json:"name"`
	Owner        string       `json:"owner"`
	BucketID     string       `json:"bucketId"`
	Payload      bytes.Buffer `json:"payload"`
}

func NewObject(b *bucket, name string) *object {
	return &object{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now().Unix(),
		LastModified: time.Now().Unix(),
		BucketID:     b.ID,
		Owner:        b.Owner,
		Name:         name,
		Version:      1,
	}
}

func (o *object) Touch() {
	o.LastModified = time.Now().Unix()
}

func (o object) TableName() string {
	return "objects"
}

func (o *object) Read(p []byte) (int, error) {
	w := bytes.NewBuffer(p)
	if _, err := o.Payload.WriteTo(w); err != nil {
		return 0, err
	}
	return len(w.Bytes()), nil
}

func (o *object) Write(p []byte) (int, error) {
	return o.Payload.Write(p)
}
