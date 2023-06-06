package osdlite

import (
	"bytes"
	"io"
	"time"

	"github.com/google/uuid"
)

type object struct {
	ID           string `json:"id"`
	CreatedAt    int64  `json:"createdAt"`
	LastModified int64  `json:"lastModified"`
	Version      int    `json:"version"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	BucketID     string `json:"bucketId"`
	Payload      []byte `json:"payload"`
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
		Payload:      make([]byte, 0, 10),
	}
}

func (o *object) Touch() {
	o.LastModified = time.Now().Unix()
}

func (o object) TableName() string {
	return "objects"
}

func (o *object) Write(p []byte) (int, error) {
	var buf bytes.Buffer
	n, err := buf.Write(p)
	if err != nil {
		return 0, err
	}
	o.Payload = buf.Bytes()
	return n, nil
}

func (o *object) Read(p []byte) (int, error) {
	if len(o.Payload) == 0 {
		return 0, io.EOF
	}
	return copy(p, o.Payload), nil
}
